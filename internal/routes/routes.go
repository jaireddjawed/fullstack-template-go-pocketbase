// Package routes declares every custom HTTP route in one place, like
// Laravel's routes/api.php. Handlers live in internal/actions.
package routes

import (
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/jaireddjawed/fullstack-template-golang/internal/actions"
)

// Register attaches all custom routes to the router. Called from the
// OnServe hook (see internal/app) and from API tests.
//
// Custom routes live under /api/app/* so they never collide with
// PocketBase's built-in /api/collections/* CRUD endpoints.
func Register(se *core.ServeEvent) {
	g := se.Router.Group("/api/app")

	g.GET("/health", actions.Health)
	g.GET("/me", actions.Me).Bind(apis.RequireAuth())
	g.GET("/posts/stats", actions.PostStats)
	g.POST("/posts/{id}/publish", actions.PublishPost).Bind(apis.RequireAuth())
}
