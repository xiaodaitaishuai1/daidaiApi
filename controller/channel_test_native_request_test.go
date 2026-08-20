package controller

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/stretchr/testify/require"
)

func TestBuildTestRequestUsesNativeAnthropicFormat(t *testing.T) {
	request := buildTestRequest("claude-test", string(constant.EndpointTypeAnthropic), nil, false)

	claudeRequest, ok := request.(*dto.ClaudeRequest)
	require.True(t, ok)
	require.Equal(t, "claude-test", claudeRequest.Model)
	require.Len(t, claudeRequest.Messages, 1)
}

func TestBuildTestRequestUsesNativeGeminiFormat(t *testing.T) {
	request := buildTestRequest("gemini-test", string(constant.EndpointTypeGemini), nil, false)

	geminiRequest, ok := request.(*dto.GeminiChatRequest)
	require.True(t, ok)
	require.Len(t, geminiRequest.Contents, 1)
	require.NotNil(t, geminiRequest.GenerationConfig.MaxOutputTokens)
	require.EqualValues(t, 3000, *geminiRequest.GenerationConfig.MaxOutputTokens)
}
