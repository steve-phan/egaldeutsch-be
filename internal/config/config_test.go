package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
  secret_key: testsecret
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
	if cfg.Jwt.SecretKey != "testsecret" {
		t.Fatalf("expected jwt secret testsecret got %s", cfg.Jwt.SecretKey)
	}
}
