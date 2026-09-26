package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSaveConfigIsPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permissions")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Simulate a config written by an older CLI version with loose perms.
	dir := filepath.Join(home, ".decision-cli")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(file, []byte("token: old\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := SaveConfig(&Config{APIURL: "https://api.example.test", Token: "secret"}); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if got := fi.Mode().Perm(); got != 0o600 {
		t.Errorf("config.yaml mode = %o, want 600", got)
	}
	di, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := di.Mode().Perm(); got != 0o700 {
		t.Errorf("config dir mode = %o, want 700", got)
	}
}
