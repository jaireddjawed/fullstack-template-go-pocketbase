package web

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewInertiaRequiresFrontendAssets(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.MkdirAll(filepath.Dir(RootTemplate), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(RootTemplate, []byte(`{{ .inertia }}{{ viteAssets "src/main.tsx" }}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := NewInertia()
	if err == nil {
		t.Fatal("expected NewInertia to reject missing Vite hot file and build manifest")
	}
	if !strings.Contains(err.Error(), "frontend assets unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewInertiaAllowsViteHotFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.MkdirAll(filepath.Dir(RootTemplate), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(RootTemplate, []byte(`{{ .inertia }}{{ viteAssets "src/main.tsx" }}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(HotFile, []byte("http://localhost:5173"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := NewInertia(); err != nil {
		t.Fatalf("NewInertia() with hot file failed: %v", err)
	}
}
