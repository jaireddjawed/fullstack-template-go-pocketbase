// fullstack-template is the project generator for this template.
//
// Interactive (TUI):
//
//	go run ./cmd/fullstack-template init
//
// Non-interactive:
//
//	go run ./cmd/fullstack-template init --no-input --name my-app \
//	  --module github.com/you/my-app --frontend next --auth pocketbase --extras shadcn
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/scaffold"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/cmd/fullstack-template/tui"
)

// defaultRepo is used when the generator is not run from inside a checkout
// of the template repository.
const defaultRepo = "https://github.com/jaireddjawed/fullstack-template-go-pocketbase.git"

func main() {
	if len(os.Args) < 2 || os.Args[1] != "init" {
		fmt.Fprintln(os.Stderr, "usage: fullstack-template init [flags]\nrun `fullstack-template init --help` for flags")
		os.Exit(2)
	}

	flags := flag.NewFlagSet("init", flag.ExitOnError)
	var (
		name     = flags.String("name", "", "project name (also the target directory)")
		module   = flags.String("module", "", "Go module path for the generated project")
		frontend = flags.String("frontend", string(scaffold.FrontendNone), "frontend: none | next | inertia")
		auth     = flags.String("auth", string(scaffold.AuthPocketBase), "auth: pocketbase | clerk | workos")
		database = flags.String("database", string(scaffold.DatabasePocketBase), "database: pocketbase")
		extras   = flags.String("extras", "", "comma-separated extras: email-verification,shadcn")
		repo     = flags.String("repo", "", "template repository to clone (defaults to the local checkout or GitHub)")
		dir      = flags.String("dir", "", "target directory (defaults to ./<name>)")
		noInput  = flags.Bool("no-input", false, "skip the TUI and use the flags as-is")
	)
	if err := flags.Parse(os.Args[2:]); err != nil {
		os.Exit(2)
	}

	cfg := scaffold.Config{
		Name:     *name,
		Module:   *module,
		Frontend: scaffold.Frontend(*frontend),
		Auth:     scaffold.Auth(*auth),
		Database: scaffold.Database(*database),
		RepoURL:  *repo,
	}
	for _, e := range strings.Split(*extras, ",") {
		if e = strings.TrimSpace(e); e != "" {
			cfg.Extras = append(cfg.Extras, scaffold.Extra(e))
		}
	}

	if !*noInput {
		collected, ok, err := tui.Run(cfg)
		if err != nil {
			fatal(err)
		}
		if !ok {
			fmt.Println("Aborted.")
			os.Exit(1)
		}
		cfg = collected
	}

	if cfg.RepoURL == "" {
		cfg.RepoURL = detectRepo()
	}
	if cfg.TargetDir = *dir; cfg.TargetDir == "" {
		cfg.TargetDir = filepath.Join(".", cfg.Name)
	}

	if err := scaffold.Generate(cfg, func(format string, args ...any) {
		fmt.Printf("• "+format+"\n", args...)
	}); err != nil {
		fatal(err)
	}

	fmt.Printf("\n✔ Generated %s\n\nNext steps:\n", cfg.TargetDir)
	for _, step := range scaffold.NextSteps(cfg) {
		fmt.Println("  " + step)
	}
}

// detectRepo prefers the local template checkout (fast, works offline) and
// falls back to the public repository.
func detectRepo() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return defaultRepo
	}
	root := strings.TrimSpace(string(out))

	gomod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil || !strings.HasPrefix(string(gomod), "module github.com/jaireddjawed/fullstack-template-go-pocketbase") {
		return defaultRepo
	}
	return root
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
