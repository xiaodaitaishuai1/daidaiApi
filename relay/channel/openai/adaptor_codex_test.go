package openai

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOpenAIResponsesSetupRequestHeaderForwardsCodexIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.141.0 (linux; x86_64)")
	c.Request.Header.Set("Originator", "codex_cli_rs")
	c.Request.Header.Set("session_id", "session-123")
	c.Request.Header.Set("X-Codex-Installation-Id", "install-123")
	c.Set("token_codex_identity_passthrough", true)

	headers := http.Header{}
	err := (&Adaptor{}).SetupRequestHeader(c, &headers, &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeResponses,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType: constant.ChannelTypeOpenAI,
			ApiKey:      "sub2-api-key",
		},
	})

	require.NoError(t, err)
	require.Equal(t, "codex_cli_rs/0.141.0 (linux; x86_64)", headers.Get("User-Agent"))
	require.Equal(t, "codex_cli_rs", headers.Get("Originator"))
	require.Equal(t, "install-123", headers.Get("X-Codex-Installation-Id"))
	require.Equal(t, "Bearer sub2-api-key", headers.Get("Authorization"))
}

func TestOpenAIResponsesSetupRequestHeaderDoesNotForwardCodexIdentityWhenTokenDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.141.0 (linux; x86_64)")
	c.Request.Header.Set("Originator", "external-client")
	c.Request.Header.Set("session_id", "session-123")
	c.Request.Header.Set("X-Codex-Installation-Id", "install-123")

	headers := http.Header{}
	err := (&Adaptor{}).SetupRequestHeader(c, &headers, &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeResponses,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType: constant.ChannelTypeOpenAI,
			ApiKey:      "sub2-api-key",
		},
	})

	require.NoError(t, err)
	require.Empty(t, headers.Get("User-Agent"))
	require.Empty(t, headers.Get("Originator"))
	require.Empty(t, headers.Get("session_id"))
	require.Empty(t, headers.Get("X-Codex-Installation-Id"))
	require.Equal(t, "Bearer sub2-api-key", headers.Get("Authorization"))
}

func TestOpenAIChatSetupRequestHeaderDoesNotForwardCodexIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.141.0 (linux; x86_64)")
	c.Request.Header.Set("X-Codex-Installation-Id", "install-123")
	c.Set("token_codex_identity_passthrough", true)

	headers := http.Header{}
	err := (&Adaptor{}).SetupRequestHeader(c, &headers, &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeChatCompletions,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType: constant.ChannelTypeOpenAI,
			ApiKey:      "sub2-api-key",
		},
	})

	require.NoError(t, err)
	require.Empty(t, headers.Get("User-Agent"))
	require.Empty(t, headers.Get("X-Codex-Installation-Id"))
}

func TestOpenAIResponsesDoRequestForwardsCodexHeadersAndBody(t *testing.T) {
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

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	c.Request.Header.Set("User-Agent", "codex_cli_rs/0.141.0 (linux; x86_64)")
	c.Request.Header.Set("Originator", "codex_cli_rs")
	c.Request.Header.Set("X-Codex-Installation-Id", "install-123")
	c.Set("token_codex_identity_passthrough", true)

	body := `{"model":"gpt-5.6-terra","client_metadata":{"x-codex-window-id":"window-123"}}`
	info := &relaycommon.RelayInfo{
		RelayMode:      relayconstant.RelayModeResponses,
		RequestURLPath: "/v1/responses",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:    constant.ChannelTypeOpenAI,
			ChannelBaseUrl: server.URL,
			ApiKey:         "sub2-api-key",
		},
	}

	resp, err := (&Adaptor{}).DoRequest(c, info, strings.NewReader(body))

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "/v1/responses", received.URL.Path)
	require.Equal(t, "codex_cli_rs/0.141.0 (linux; x86_64)", received.Header.Get("User-Agent"))
	require.Equal(t, "codex_cli_rs", received.Header.Get("Originator"))
	require.Equal(t, "install-123", received.Header.Get("X-Codex-Installation-Id"))
	require.Equal(t, "Bearer sub2-api-key", received.Header.Get("Authorization"))
	require.Equal(t, body, receivedBody)
}
