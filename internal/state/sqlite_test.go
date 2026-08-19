package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenSQLiteCreatesPrivateDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state.db")
	store, err := OpenSQLite(path)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer store.Close()

	if _, err := store.DB().Exec("CREATE TABLE check_state (id INTEGER PRIMARY KEY)"); err != nil {
		t.Fatalf("create table: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat database: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("database permissions = %o, want 600", got)
	}
}

func TestOpenSQLiteRequiresPath(t *testing.T) {
	if _, err := OpenSQLite("  "); err == nil {
		t.Fatal("expected empty path to fail")
	}
}
