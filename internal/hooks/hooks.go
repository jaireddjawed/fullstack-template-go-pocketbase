// Package hooks centralizes PocketBase event hooks — logic that runs when
// records change, regardless of whether the change came from the REST API,
// the admin dashboard, or Go code. Comparable to Laravel model observers.
package hooks

import (
	"github.com/pocketbase/pocketbase/core"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/models"
	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/services"
)

// Register binds all event hooks to the app.
func Register(app core.App) {
	// Auto-generate a slug from the title whenever a post is created
	// without one.
	app.OnRecordCreate("posts").BindFunc(func(e *core.RecordEvent) error {
		post := models.NewPost(e.Record)
		data := post.Data()
		if data.Slug == "" {
			data.Slug = services.Slugify(data.Title)
			post.Apply(data)
		}
		return e.Next()
	})
}
