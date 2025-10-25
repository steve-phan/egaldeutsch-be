package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgpkg "egaldeutsch-be/internal/config"
)

func TestLoadConfigFromFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "cfgtest")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	content := `server:
  port: "8080"
  host: "localhost"
database:
  host: "localhost"
  port: 5432
  user: postgres
  password: postgres
  dbname: egaldeutsch
  sslmode: disable
jwt:
  secret_key: this-is-a-very-secure-secret-key-with-32-plus-characters
  issuer: egaldeutsch
  expiration_hours: 24
  refresh_token_expiration_days: 30
`
	fpath := filepath.Join(dir, "config.yaml")
	if err := ioutil.WriteFile(fpath, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg, err := cfgpkg.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Server.Port != "8080" {
		t.Fatalf("expected port 8080 got %s", cfg.Server.Port)
	}
	if cfg.Jwt.SecretKey != "this-is-a-very-secure-secret-key-with-32-plus-characters" {
		t.Fatalf("expected jwt secret this-is-a-very-secure-secret-key-with-32-plus-characters got %s", cfg.Jwt.SecretKey)
	}
}

func TestJwtConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    cfgpkg.JwtConfig
		shouldErr bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
				Issuer:                     "egaldeutsch",
				ExpirationHours:            24,
				RefreshTokenExpirationDays: 30,
			},
			shouldErr: false,
		},
		{
			name: "empty secret key",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "",
				Issuer:                     "egaldeutsch",
				ExpirationHours:            24,
				RefreshTokenExpirationDays: 30,
			},
			shouldErr: true,
			errorMsg:  "JWT secret key cannot be empty",
		},
		{
			name: "short secret key",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "tooshort",
				Issuer:                     "egaldeutsch",
				ExpirationHours:            24,
				RefreshTokenExpirationDays: 30,
			},
			shouldErr: true,
			errorMsg:  "JWT secret key must be at least 32 characters",
		},
		{
			name: "empty issuer",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
				Issuer:                     "",
				ExpirationHours:            24,
				RefreshTokenExpirationDays: 30,
			},
			shouldErr: true,
			errorMsg:  "JWT issuer cannot be empty",
		},
		{
			name: "zero expiration hours",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
				Issuer:                     "egaldeutsch",
				ExpirationHours:            0,
				RefreshTokenExpirationDays: 30,
			},
			shouldErr: true,
			errorMsg:  "JWT expiration hours must be positive",
		},
		{
			name: "too long expiration hours",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
				Issuer:                     "egaldeutsch",
				ExpirationHours:            10000, // > 1 year
				RefreshTokenExpirationDays: 30,
			},
			shouldErr: true,
			errorMsg:  "JWT expiration hours too long",
		},
		{
			name: "zero refresh token expiration days",
			config: cfgpkg.JwtConfig{
				SecretKey:                  "this-is-a-very-secure-secret-key-with-32-plus-characters",
				Issuer:                     "egaldeutsch",
				ExpirationHours:            24,
				RefreshTokenExpirationDays: 0,
			},
			shouldErr: true,
			errorMsg:  "JWT refresh token expiration days must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error containing '%s', got nil", tt.errorMsg)
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Fatalf("expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}
