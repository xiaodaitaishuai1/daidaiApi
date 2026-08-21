package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAppendCodexIdentityTraceAdminInfoAddsSafeTrace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	common.SetContextKey(c, constant.ContextKeyCodexIdentityTrace, "originator=sha1:123456789abc")
	adminInfo := map[string]interface{}{}

	AppendCodexIdentityTraceAdminInfo(c, adminInfo)

	require.Equal(t, "originator=sha1:123456789abc", adminInfo["codex_identity_trace"])
}

func TestAppendCodexIdentityTraceAdminInfoSkipsMissingTrace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	adminInfo := map[string]interface{}{}

	AppendCodexIdentityTraceAdminInfo(c, adminInfo)

	require.Empty(t, adminInfo)
}
