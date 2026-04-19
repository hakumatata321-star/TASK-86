package config_test

import (
	"os"
	"testing"

	"w2t86/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DB_PATH")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("TIMEZONE")
	os.Unsetenv("ENCRYPTION_KEY")
	os.Unsetenv("SESSION_SECRET")
	os.Unsetenv("BANNED_WORDS")

	cfg := config.Load()

	if cfg.Port != "3000" {
		t.Errorf("expected default Port=3000, got %q", cfg.Port)
	}
	if cfg.DBPath != "data/portal.db" {
		t.Errorf("expected default DBPath=data/portal.db, got %q", cfg.DBPath)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("expected default AppEnv=development, got %q", cfg.AppEnv)
	}
	if cfg.Timezone != "UTC" {
		t.Errorf("expected default Timezone=UTC, got %q", cfg.Timezone)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DB_PATH", "/tmp/test.db")
	t.Setenv("APP_ENV", "production")
	t.Setenv("TIMEZONE", "America/New_York")
	t.Setenv("ENCRYPTION_KEY", "aabbcc")
	t.Setenv("SESSION_SECRET", "mysecret")
	t.Setenv("BANNED_WORDS", "spam,abuse")

	cfg := config.Load()

	if cfg.Port != "8080" {
		t.Errorf("expected Port=8080, got %q", cfg.Port)
	}
	if cfg.DBPath != "/tmp/test.db" {
		t.Errorf("expected DBPath=/tmp/test.db, got %q", cfg.DBPath)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("expected AppEnv=production, got %q", cfg.AppEnv)
	}
	if cfg.Timezone != "America/New_York" {
		t.Errorf("expected Timezone=America/New_York, got %q", cfg.Timezone)
	}
	if cfg.EncryptionKey != "aabbcc" {
		t.Errorf("expected EncryptionKey=aabbcc, got %q", cfg.EncryptionKey)
	}
	if cfg.SessionSecret != "mysecret" {
		t.Errorf("expected SessionSecret=mysecret, got %q", cfg.SessionSecret)
	}
	if cfg.BannedWords != "spam,abuse" {
		t.Errorf("expected BannedWords=spam,abuse, got %q", cfg.BannedWords)
	}
}

func TestValidate_MissingEncryptionKey(t *testing.T) {
	cfg := &config.Config{
		SessionSecret: "some-secret",
		EncryptionKey: "",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty EncryptionKey, got nil")
	}
}

func TestValidate_MissingSessionSecret(t *testing.T) {
	cfg := &config.Config{
		EncryptionKey: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		SessionSecret: "",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty SessionSecret, got nil")
	}
}

func TestValidate_InvalidHexKey(t *testing.T) {
	cfg := &config.Config{
		EncryptionKey: "not-valid-hex!!!",
		SessionSecret: "some-secret",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for non-hex EncryptionKey, got nil")
	}
}

func TestValidate_ShortKey(t *testing.T) {
	// 16 bytes = 32 hex chars, but we need 32 bytes = 64 hex chars.
	cfg := &config.Config{
		EncryptionKey: "aabbccddeeff00112233445566778899", // only 32 chars
		SessionSecret: "some-secret",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for 16-byte key, got nil")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		EncryptionKey: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		SessionSecret: "a-long-random-session-secret-value",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}

func TestValidate_BothMissing(t *testing.T) {
	cfg := &config.Config{}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error when both secrets are missing, got nil")
	}
	msg := err.Error()
	if len(msg) == 0 {
		t.Error("expected non-empty error message")
	}
}
