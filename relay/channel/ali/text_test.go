package ali

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	"github.com/stretchr/testify/require"
)

func TestRequestOpenAI2AliDoesNotInjectOmittedTopP(t *testing.T) {
	converted := requestOpenAI2Ali(dto.GeneralOpenAIRequest{Model: "qwen-plus"})

	require.Nil(t, converted.TopP)
}

func TestRequestOpenAI2AliClampsExplicitTopP(t *testing.T) {
	topP := 1.0
	converted := requestOpenAI2Ali(dto.GeneralOpenAIRequest{Model: "qwen-plus", TopP: &topP})

	require.NotNil(t, converted.TopP)
	require.Equal(t, 0.99, *converted.TopP)
}
