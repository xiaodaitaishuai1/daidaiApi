package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestHardDeletePurgesAuthenticationData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()
	if err := db.AutoMigrate(&User{}, &Token{}, &TwoFA{}, &TwoFABackupCode{}, &PasskeyCredential{}, &UserOAuthBinding{}); err != nil {
		t.Fatal(err)
	}
	user := &User{Username: "hard-delete-test", Password: "x", Status: 1}
	if err := db.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&Token{UserId: user.Id, Key: "hard-delete-token"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&TwoFA{UserId: user.Id, Secret: "secret"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&PasskeyCredential{UserID: user.Id, CredentialID: "credential", PublicKey: "key"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&UserOAuthBinding{UserId: user.Id, ProviderId: 1, ProviderUserId: "oauth-user"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := user.HardDelete(); err != nil {
		t.Fatal(err)
	}
	for _, query := range []struct {
		name  string
		model any
		where string
	}{
		{"user", &User{}, "id = ?"}, {"token", &Token{}, "user_id = ?"}, {"twofa", &TwoFA{}, "user_id = ?"},
		{"passkey", &PasskeyCredential{}, "user_id = ?"}, {"oauth", &UserOAuthBinding{}, "user_id = ?"},
	} {
		var count int64
		if err := db.Unscoped().Model(query.model).Where(query.where, user.Id).Count(&count).Error; err != nil {
			t.Fatal(query.name, err)
		}
		if count != 0 {
			t.Fatalf("%s data remains after hard delete: %d", query.name, count)
		}
	}
}
