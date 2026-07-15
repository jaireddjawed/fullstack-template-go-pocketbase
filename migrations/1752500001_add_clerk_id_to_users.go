package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Adds the clerk_id field used to link PocketBase user records to Clerk
// identities. Unique for non-empty values (users created outside Clerk,
// e.g. by the seeder, have no clerk_id).
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		users.Fields.Add(&core.TextField{Name: "clerk_id", Max: 100})
		users.AddIndex("idx_users_clerk_id", true, "clerk_id", "clerk_id != ''")

		return app.Save(users)
	}, func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		users.RemoveIndex("idx_users_clerk_id")
		users.Fields.RemoveByName("clerk_id")

		return app.Save(users)
	})
}
