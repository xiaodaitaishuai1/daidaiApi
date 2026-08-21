package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTokenCodexIdentityPassthroughPersistsOnCreateAndUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()
	require.NoError(t, db.AutoMigrate(&Token{}))

	token := &Token{
		UserId:                   1,
		Key:                      "codex-identity-test-key",
		Name:                     "codex identity test",
		CreatedTime:              1,
		AccessedTime:             1,
		ExpiredTime:              -1,
		CodexIdentityPassthrough: true,
	}
	require.NoError(t, token.Insert())

	var stored Token
	require.NoError(t, db.First(&stored, token.Id).Error)
	require.True(t, stored.CodexIdentityPassthrough)

	stored.CodexIdentityPassthrough = false
	require.NoError(t, stored.Update())

	var updated Token
	require.NoError(t, db.First(&updated, token.Id).Error)
	require.False(t, updated.CodexIdentityPassthrough)
}
