package middleware

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
)

func TestDecompressRequestMiddlewareSupportsZstd(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var compressed bytes.Buffer
	encoder, err := zstd.NewWriter(&compressed)
	require.NoError(t, err)
	_, err = encoder.Write([]byte(`{"message":"hello"}`))
	require.NoError(t, err)
	require.NoError(t, encoder.Close())

	router := gin.New()
	router.Use(DecompressRequestMiddleware())
	router.POST("/", func(c *gin.Context) {
		body, readErr := io.ReadAll(c.Request.Body)
		require.NoError(t, readErr)
		c.Data(200, "application/json", body)
	})

	request := httptest.NewRequest("POST", "/", &compressed)
	request.Header.Set("Content-Encoding", "zstd")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	require.Equal(t, 200, response.Code)
	require.JSONEq(t, `{"message":"hello"}`, response.Body.String())
}
