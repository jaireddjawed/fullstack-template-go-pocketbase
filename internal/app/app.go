// Package app wires the PocketBase application together: CLI commands,
// event hooks, and custom routes. Everything that main.go boots and that
// tests need to replicate lives behind Bind().
package app

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/commands"
	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/hooks"
	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/routes"
)

// New builds the full PocketBase application used by main.go.
func New() *pocketbase.PocketBase {
	pb := pocketbase.New()

	// `migrate` CLI command + automigrate: editing collections in the admin
	// dashboard writes a Go migration file into ./migrations.
	migratecmd.MustRegister(pb, pb.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangGo,
		Automigrate:  true,
		Dir:          "migrations",
	})

	commands.Register(pb)
	Bind(pb)

	return pb
}

// Bind attaches event hooks and route registration to any core.App.
// Tests call this on a tests.TestApp so scenarios exercise the exact
// same hooks and routes as production.
func Bind(app core.App) {
	hooks.Register(app)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		routes.Register(se)
		return se.Next()
	})
}
