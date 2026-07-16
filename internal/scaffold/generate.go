package scaffold

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Logf receives human-readable progress lines during generation.
type Logf func(format string, args ...any)

// Generate creates a new project in cfg.TargetDir from the template branch
// selected by cfg. Steps: clone → strip history → rewrite module path →
// write Docker files → apply extras → fresh git init.
func Generate(cfg Config, logf Logf) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	branch, err := cfg.Branch()
	if err != nil {
		return err
	}

	if err := ensureEmptyDir(cfg.TargetDir); err != nil {
		return err
	}

	logf("Cloning branch %s from %s", branch, cfg.RepoURL)
	if out, err := exec.Command("git", "clone", "--depth", "1", "--branch", branch, cfg.RepoURL, cfg.TargetDir).CombinedOutput(); err != nil {
		return fmt.Errorf("git clone: %w\n%s", err, out)
	}
	if err := os.RemoveAll(filepath.Join(cfg.TargetDir, ".git")); err != nil {
		return err
	}

	oldModule, err := modulePath(filepath.Join(cfg.TargetDir, "go.mod"))
	if err != nil {
		return err
	}

	logf("Rewriting module %s -> %s", oldModule, cfg.Module)
	if err := RewriteModule(cfg.TargetDir, oldModule, cfg.Module); err != nil {
		return err
	}

	logf("Writing Docker files")
	if err := writeDockerFiles(cfg); err != nil {
		return fmt.Errorf("docker files: %w", err)
	}

	for _, extra := range cfg.Extras {
		logf("Adding extra: %s", extra)
		if err := applyExtra(cfg, extra); err != nil {
			return fmt.Errorf("extra %s: %w", extra, err)
		}
	}

	logf("Initializing fresh git repository")
	for _, args := range [][]string{
		{"init", "-q"},
		{"add", "-A"},
		// -c identity fallbacks keep the commit working where git has no
		// global identity configured (fresh machines, CI).
		{"-c", "user.name=fullstack-template", "-c", "user.email=generator@localhost", "commit", "-q", "-m", "Initial commit from fullstack-template"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = cfg.TargetDir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, out)
		}
	}

	return nil
}

func ensureEmptyDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("target directory %s is not empty", dir)
	}
	// git clone wants to create the directory itself
	return os.Remove(dir)
}

// modulePath reads the module directive from a go.mod file.
func modulePath(gomod string) (string, error) {
	f, err := os.Open(gomod)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if module, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(module), nil
		}
	}
	return "", fmt.Errorf("no module directive in %s", gomod)
}

// rewriteExts are the file types that may reference the Go module path
// (imports, docs, configs like tygo.yaml).
var rewriteExts = []string{".go", ".mod", ".md", ".yaml", ".yml", ".ts", ".tsx", ".json"}

var skipDirs = map[string]bool{".git": true, "node_modules": true, "pb_data": true, "dist": true, ".next": true}

// RewriteModule replaces the Go module path in every relevant file under dir.
func RewriteModule(dir, oldModule, newModule string) error {
	if oldModule == newModule {
		return nil
	}

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !slices.Contains(rewriteExts, filepath.Ext(path)) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Contains(content, []byte(oldModule)) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		updated := bytes.ReplaceAll(content, []byte(oldModule), []byte(newModule))
		return os.WriteFile(path, updated, info.Mode())
	})
}

// NextSteps returns the checklist printed after a successful generation.
func NextSteps(cfg Config) []string {
	steps := []string{
		fmt.Sprintf("cd %s", cfg.Name),
		"make dev                # backend on http://127.0.0.1:8090",
		"make seed               # demo data (demo@example.com / password123)",
	}

	if cfg.Frontend != FrontendNone {
		steps = append(steps, "cd frontend && npm install && npm run dev")
	}

	if cfg.Auth == AuthClerk {
		steps = append(steps,
			"cp frontend/.env.example frontend/.env.local   # add your Clerk keys",
			"export CLERK_SECRET_KEY=sk_...                 # before `make dev`",
		)
	}

	steps = append(steps, "docker compose up --build      # containerized run")

	if slices.Contains(cfg.Extras, ExtraShadcn) {
		steps = append(steps,
			"cd frontend && npx shadcn@latest init          # set up shadcn/ui",
			"npx shadcn@latest add button                   # then add components",
		)
	}

	steps = append(steps, "See docs/ for architecture, database, types, and testing guides.")
	return steps
}
