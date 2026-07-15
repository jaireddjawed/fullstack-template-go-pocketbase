package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Creates the demo "posts" collection.
//
// Every migration has an up and a down function and runs inside a
// transaction. New migrations are created with:
//
//	go run . migrate create "some_change"   (or: make migration name=some_change)
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		posts := core.NewBaseCollection("posts")

		posts.Fields.Add(
			&core.TextField{Name: "title", Required: true, Max: 200},
			&core.TextField{Name: "slug", Max: 250},
			&core.EditorField{Name: "content"},
			&core.BoolField{Name: "published"},
			&core.RelationField{
				Name:          "owner",
				Required:      true,
				CollectionId:  users.Id,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)

		posts.AddIndex("idx_posts_slug", false, "slug", "")
		posts.AddIndex("idx_posts_owner", false, "owner", "")

		// API rules (PocketBase's row-level authorization):
		// anyone can read published posts, owners can do everything with their own.
		posts.ListRule = types.Pointer("published = true || owner = @request.auth.id")
		posts.ViewRule = types.Pointer("published = true || owner = @request.auth.id")
		posts.CreateRule = types.Pointer("@request.auth.id != '' && owner = @request.auth.id")
		posts.UpdateRule = types.Pointer("owner = @request.auth.id")
		posts.DeleteRule = types.Pointer("owner = @request.auth.id")

		return app.Save(posts)
	}, func(app core.App) error {
		posts, err := app.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}
		return app.Delete(posts)
	})
}
