package codex

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newCodexHeaderTestContext(headers map[string]string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{"model":"gpt-5-codex"}`))
	for name, value := range headers {
		c.Request.Header.Set(name, value)
	}
	return c
}

func newCodexHeaderTestInfo(baseURL string, relayMode int) *relaycommon.RelayInfo {
	return &relaycommon.RelayInfo{
		RelayMode: relayMode,
		IsStream:  relayMode == relayconstant.RelayModeResponses,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:    constant.ChannelTypeCodex,
			ChannelBaseUrl: baseURL,
			ApiKey:         `{"access_token":"upstream-token","account_id":"acct-123"}`,
			HeadersOverride: map[string]interface{}{
				"User-Agent": "codex_cli_rs/9.9.9 (override)",
			},
		},
	}
}

func TestSetupRequestHeaderForwardsCodexIdentityHeaders(t *testing.T) {
	c := newCodexHeaderTestContext(map[string]string{
		"User-Agent":              "codex_cli_rs/0.141.0 (linux; x86_64)",
		"Originator":              "codex_cli_rs",
		"session_id":              "session-123",
		"X-Codex-Installation-Id": "install-123",
		"X-Codex-Window-Id":       "window-123",
		"X-Codex-Turn-State":      "state-123",
		"X-Codex-Turn-Metadata":   `{"turn_id":"turn-123"}`,
		"Authorization":           "Bearer client-token",
		"Cookie":                  "session=secret",
		"Accept-Encoding":         "gzip",
	})
	c.Request.Header["x-codex-raw-fingerprint"] = []string{"raw-fingerprint"}
	info := newCodexHeaderTestInfo("https://sub2.example.com", relayconstant.RelayModeResponses)
	headers := http.Header{}

	err := (&Adaptor{}).SetupRequestHeader(c, &headers, info)

	require.NoError(t, err)
	require.Equal(t, "codex_cli_rs/0.141.0 (linux; x86_64)", headers.Get("User-Agent"))
	require.Equal(t, "codex_cli_rs", headers.Get("Originator"))
	require.Equal(t, "session-123", headers.Get("Session_id"))
	require.Equal(t, "install-123", headers.Get("X-Codex-Installation-Id"))
	require.Equal(t, "window-123", headers.Get("X-Codex-Window-Id"))
	require.Equal(t, "state-123", headers.Get("X-Codex-Turn-State"))
	require.Equal(t, `{"turn_id":"turn-123"}`, headers.Get("X-Codex-Turn-Metadata"))
	require.Equal(t, "raw-fingerprint", headers.Get("X-Codex-Raw-Fingerprint"))
	require.Equal(t, "Bearer upstream-token", headers.Get("Authorization"))
	require.Equal(t, "acct-123", headers.Get("chatgpt-account-id"))
	require.Empty(t, headers.Get("Cookie"))
	require.Empty(t, headers.Get("Accept-Encoding"))
	require.Empty(t, headers.Get("Host"))
}

func TestSetupRequestHeaderDoesNotSynthesizeMissingUserAgent(t *testing.T) {
	c := newCodexHeaderTestContext(map[string]string{"X-Codex-Installation-Id": "install-123"})
	info := newCodexHeaderTestInfo("https://sub2.example.com", relayconstant.RelayModeResponses)
	headers := http.Header{}

	err := (&Adaptor{}).SetupRequestHeader(c, &headers, info)

	require.NoError(t, err)
	require.Empty(t, headers.Get("User-Agent"))
	require.Equal(t, "codex_cli_rs", headers.Get("Originator"))
}

func TestConvertOpenAIResponsesRequestPreservesClientMetadata(t *testing.T) {
	var request dto.OpenAIResponsesRequest
	err := common.Unmarshal([]byte(`{"model":"gpt-5-codex","input":"hi","client_metadata":{"x-codex-window-id":"window-123"}}`), &request)
	require.NoError(t, err)

	converted, err := (&Adaptor{}).ConvertOpenAIResponsesRequest(nil, &relaycommon.RelayInfo{
		RelayMode:   relayconstant.RelayModeResponses,
		ChannelMeta: &relaycommon.ChannelMeta{},
	}, request)
	require.NoError(t, err)

	body, err := common.Marshal(converted)
	require.NoError(t, err)
	var payload map[string]interface{}
	require.NoError(t, common.Unmarshal(body, &payload))
	metadata, ok := payload["client_metadata"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "window-123", metadata["x-codex-window-id"])
}

func TestDoApiRequestUsesHeaderOverrideAndPreservesResponsesPaths(t *testing.T) {
	tests := []struct {
		name string
		mode int
		path string
	}{
		{name: "responses", mode: relayconstant.RelayModeResponses, path: "/backend-api/codex/responses"},
		{name: "compact", mode: relayconstant.RelayModeResponsesCompact, path: "/backend-api/codex/responses/compact"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.InitHttpClient()
			var received *http.Request
			var receivedBody string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				received = r
				body, _ := io.ReadAll(r.Body)
				receivedBody = string(body)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"resp-1"}`))
			}))
			defer server.Close()

			c := newCodexHeaderTestContext(map[string]string{
				"User-Agent":              "codex_cli_rs/0.141.0 (linux; x86_64)",
				"Originator":              "codex_cli_rs",
				"X-Codex-Installation-Id": "install-123",
			})
			info := newCodexHeaderTestInfo(server.URL, tt.mode)
			body := `{"model":"gpt-5-codex","client_metadata":{"x-codex-window-id":"window-123"}}`
			resp, err := (&Adaptor{}).DoRequest(c, info, strings.NewReader(body))

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tt.path, received.URL.Path)
			require.Equal(t, "codex_cli_rs/9.9.9 (override)", received.Header.Get("User-Agent"))
			require.Equal(t, "install-123", received.Header.Get("X-Codex-Installation-Id"))
			require.Equal(t, "Bearer upstream-token", received.Header.Get("Authorization"))
			require.Empty(t, received.Header.Get("Cookie"))
			require.Equal(t, body, receivedBody)
		})
	}
}
