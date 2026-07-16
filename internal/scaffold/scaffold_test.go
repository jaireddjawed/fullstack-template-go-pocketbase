package scaffold_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaireddjawed/fullstack-template-golang/internal/scaffold"
)

func TestBranchMatrix(t *testing.T) {
	cases := []struct {
		frontend scaffold.Frontend
		auth     scaffold.Auth
		branch   string
		wantErr  bool
	}{
		{scaffold.FrontendNone, scaffold.AuthPocketBase, "main", false},
		{scaffold.FrontendNext, scaffold.AuthPocketBase, "nextjs", false},
		{scaffold.FrontendNext, scaffold.AuthClerk, "nextjs-clerk", false},
		{scaffold.FrontendInertia, scaffold.AuthPocketBase, "react-inertia", false},
		{scaffold.FrontendInertia, scaffold.AuthClerk, "", true},
		{scaffold.FrontendNone, scaffold.AuthClerk, "", true},
		{scaffold.FrontendNext, scaffold.AuthWorkOS, "", true},
	}

	for _, tc := range cases {
		cfg := scaffold.Config{Frontend: tc.frontend, Auth: tc.auth, Database: scaffold.DatabasePocketBase}
		branch, err := cfg.Branch()
		if tc.wantErr != (err != nil) {
			t.Errorf("Branch(%s,%s) error = %v, wantErr %v", tc.frontend, tc.auth, err, tc.wantErr)
		}
		if branch != tc.branch {
			t.Errorf("Branch(%s,%s) = %q, want %q", tc.frontend, tc.auth, branch, tc.branch)
		}
	}
}

func TestPostgresNotImplemented(t *testing.T) {
	cfg := scaffold.Config{
		Frontend: scaffold.FrontendNext,
		Auth:     scaffold.AuthPocketBase,
		Database: scaffold.DatabasePostgres,
	}
	if _, err := cfg.Branch(); err == nil {
		t.Error("expected an error for the postgres database option")
	}
}

func TestValidate(t *testing.T) {
	valid := scaffold.Config{
		Name:     "my-app",
		Module:   "github.com/acme/my-app",
		Frontend: scaffold.FrontendNext,
		Auth:     scaffold.AuthPocketBase,
		Database: scaffold.DatabasePocketBase,
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid config rejected: %v", err)
	}

	badName := valid
	badName.Name = "my app!"
	if err := badName.Validate(); err == nil {
		t.Error("expected error for invalid project name")
	}

	badModule := valid
	badModule.Module = "not a module"
	if err := badModule.Validate(); err == nil {
		t.Error("expected error for invalid module path")
	}

	badExtra := valid
	badExtra.Extras = []scaffold.Extra{scaffold.ExtraEmailVerification}
	badExtra.Auth = scaffold.AuthClerk
	if err := badExtra.Validate(); err == nil {
		t.Error("expected error: email-verification is not applicable with clerk auth")
	}
}

func TestChoiceAvailability(t *testing.T) {
	for _, c := range scaffold.AuthChoices(scaffold.FrontendInertia) {
		if c.Value == string(scaffold.AuthClerk) && !c.Disabled {
			t.Error("clerk should be disabled for the inertia frontend")
		}
	}
	for _, c := range scaffold.ExtraChoices(scaffold.FrontendNone, scaffold.AuthPocketBase) {
		if c.Value == string(scaffold.ExtraShadcn) && !c.Disabled {
			t.Error("shadcn should be disabled without a frontend")
		}
	}
	for _, c := range scaffold.DatabaseChoices() {
		if c.Value == string(scaffold.DatabasePostgres) && !c.Disabled {
			t.Error("postgres should be marked as disabled")
		}
	}
}

func TestRewriteModule(t *testing.T) {
	dir := t.TempDir()
	old, new := "github.com/old/mod", "github.com/new/mod"

	write := func(rel, content string) {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("go.mod", "module "+old+"\n")
	write("main.go", `import "`+old+`/internal/app"`)
	write("node_modules/pkg/index.ts", old) // must be skipped
	write("script.sh", old)                 // unknown extension, skipped

	if err := scaffold.RewriteModule(dir, old, new); err != nil {
		t.Fatal(err)
	}

	assertContains := func(rel, want string) {
		content, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), want) {
			t.Errorf("%s: expected to contain %q, got %q", rel, want, content)
		}
	}

	assertContains("go.mod", new)
	assertContains("main.go", new)
	assertContains("node_modules/pkg/index.ts", old)
	assertContains("script.sh", old)
}

// TestGenerateEndToEnd clones the local repository's main branch into a
// temp dir and checks the result. Skipped when the local main branch is
// unavailable (e.g. shallow CI checkouts).
func TestGenerateEndToEnd(t *testing.T) {
	repoRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Skip("not in a git repository")
	}
	repo := strings.TrimSpace(string(repoRoot))

	if err := exec.Command("git", "-C", repo, "rev-parse", "--verify", "main").Run(); err != nil {
		t.Skip("local main branch not available")
	}

	target := filepath.Join(t.TempDir(), "generated")
	cfg := scaffold.Config{
		Name:      "my-app",
		Module:    "github.com/acme/my-app",
		Frontend:  scaffold.FrontendNone,
		Auth:      scaffold.AuthPocketBase,
		Database:  scaffold.DatabasePocketBase,
		Extras:    []scaffold.Extra{scaffold.ExtraDocker, scaffold.ExtraEmailVerification},
		RepoURL:   repo,
		TargetDir: target,
	}

	if err := scaffold.Generate(cfg, t.Logf); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	gomod, err := os.ReadFile(filepath.Join(target, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(gomod), "module github.com/acme/my-app") {
		t.Errorf("go.mod module not rewritten: %s", gomod)
	}

	for _, rel := range []string{"Dockerfile", "docker-compose.yml", ".dockerignore"} {
		if _, err := os.Stat(filepath.Join(target, rel)); err != nil {
			t.Errorf("expected %s to be generated: %v", rel, err)
		}
	}

	migrations, err := filepath.Glob(filepath.Join(target, "migrations", "*_require_verified_email.go"))
	if err != nil || len(migrations) != 1 {
		t.Errorf("expected the email verification migration, found %v", migrations)
	}

	// Fresh history: exactly one commit.
	out, err := exec.Command("git", "-C", target, "rev-list", "--count", "HEAD").Output()
	if err != nil || strings.TrimSpace(string(out)) != "1" {
		t.Errorf("expected exactly 1 commit in the generated repo, got %q (%v)", out, err)
	}
}
