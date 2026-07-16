// Package scaffold generates new projects from this template. It contains
// all generator logic — the option matrix, cloning, module rewriting, and
// extras — kept free of TUI concerns so it can be unit tested and driven
// non-interactively (cmd/fullstack-template).
package scaffold

import (
	"fmt"
	"regexp"
	"slices"
)

type Frontend string

const (
	FrontendNone    Frontend = "none"
	FrontendNext    Frontend = "next"
	FrontendInertia Frontend = "inertia"
)

type Auth string

const (
	AuthPocketBase Auth = "pocketbase"
	AuthClerk      Auth = "clerk"
	AuthWorkOS     Auth = "workos"
)

type Database string

const (
	DatabasePocketBase Database = "pocketbase"
	DatabasePostgres   Database = "postgres"
)

type Extra string

const (
	ExtraDocker            Extra = "docker"
	ExtraEmailVerification Extra = "email-verification"
	ExtraShadcn            Extra = "shadcn"
)

// Config is a fully specified generation request.
type Config struct {
	Name     string // project directory / display name
	Module   string // Go module path of the generated project
	Frontend Frontend
	Auth     Auth
	Database Database
	Extras   []Extra

	RepoURL   string // template repository (local path or URL) to clone from
	TargetDir string // where to generate; must not exist or be empty
}

// branchMatrix maps supported frontend+auth combinations to template branches.
var branchMatrix = map[Frontend]map[Auth]string{
	FrontendNone:    {AuthPocketBase: "main"},
	FrontendNext:    {AuthPocketBase: "nextjs", AuthClerk: "nextjs-clerk"},
	FrontendInertia: {AuthPocketBase: "react-inertia"},
}

// Branch resolves the template branch for the chosen stack.
func (c Config) Branch() (string, error) {
	if c.Database == DatabasePostgres {
		return "", fmt.Errorf("database %q is not implemented yet", c.Database)
	}
	if branch, ok := branchMatrix[c.Frontend][c.Auth]; ok {
		return branch, nil
	}
	return "", fmt.Errorf("no template branch implements frontend=%s auth=%s", c.Frontend, c.Auth)
}

var moduleRe = regexp.MustCompile(`^[a-zA-Z0-9._~\-/]+$`)
var nameRe = regexp.MustCompile(`^[a-zA-Z0-9._\-]+$`)

// Validate checks the config before generation.
func (c Config) Validate() error {
	if !nameRe.MatchString(c.Name) {
		return fmt.Errorf("invalid project name %q (letters, digits, . _ - only)", c.Name)
	}
	if !moduleRe.MatchString(c.Module) {
		return fmt.Errorf("invalid Go module path %q", c.Module)
	}
	if _, err := c.Branch(); err != nil {
		return err
	}
	for _, extra := range c.Extras {
		if !slices.ContainsFunc(ExtraChoices(c.Frontend, c.Auth), func(o Choice) bool {
			return o.Value == string(extra) && !o.Disabled
		}) {
			return fmt.Errorf("extra %q is not available for frontend=%s auth=%s", extra, c.Frontend, c.Auth)
		}
	}
	return nil
}

// Choice is a selectable option presented by the TUI (or --help output).
type Choice struct {
	Value    string
	Label    string
	Detail   string
	Disabled bool // shown, but not selectable ("coming soon" / not applicable)
}

func FrontendChoices() []Choice {
	return []Choice{
		{Value: string(FrontendNone), Label: "None", Detail: "backend only (branch: main)"},
		{Value: string(FrontendNext), Label: "Next.js", Detail: "separate Next.js app (branch: nextjs*)"},
		{Value: string(FrontendInertia), Label: "React + Inertia", Detail: "monolith, Go serves React (branch: react-inertia)"},
	}
}

func AuthChoices(f Frontend) []Choice {
	return []Choice{
		{Value: string(AuthPocketBase), Label: "PocketBase", Detail: "built-in email/password auth"},
		{
			Value:    string(AuthClerk),
			Label:    "Clerk",
			Detail:   "hosted auth mapped to PocketBase users",
			Disabled: f != FrontendNext, // only the Next.js stack has a Clerk integration
		},
		{Value: string(AuthWorkOS), Label: "WorkOS", Detail: "coming soon", Disabled: true},
	}
}

func DatabaseChoices() []Choice {
	return []Choice{
		{Value: string(DatabasePocketBase), Label: "PocketBase (SQLite)", Detail: "embedded, zero-config"},
		{Value: string(DatabasePostgres), Label: "Postgres", Detail: "coming soon", Disabled: true},
	}
}

func ExtraChoices(f Frontend, a Auth) []Choice {
	return []Choice{
		{Value: string(ExtraDocker), Label: "Docker", Detail: "Dockerfile + docker-compose.yml"},
		{
			Value:    string(ExtraEmailVerification),
			Label:    "Email verification",
			Detail:   "require a verified email to log in (migration)",
			Disabled: a != AuthPocketBase, // Clerk manages verification itself
		},
		{
			Value:    string(ExtraShadcn),
			Label:    "shadcn/ui",
			Detail:   "prints the shadcn init steps after generation",
			Disabled: f == FrontendNone,
		},
	}
}
