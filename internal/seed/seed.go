// Package seed populates the database with development data,
// like Laravel's database seeders. Run with: go run . seed
package seed

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/models"
)

const (
	DemoEmail    = "demo@example.com"
	DemoPassword = "password123"
)

// Run seeds the database. It is idempotent: records are looked up before
// being created, so running it twice does not duplicate data.
func Run(app core.App) error {
	user, err := ensureDemoUser(app)
	if err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	if err := ensureDemoPosts(app, user); err != nil {
		return fmt.Errorf("seed posts: %w", err)
	}

	app.Logger().Info("database seeded", "demo user", DemoEmail)
	fmt.Printf("Seeded. Demo login: %s / %s\n", DemoEmail, DemoPassword)
	return nil
}

func ensureDemoUser(app core.App) (*models.User, error) {
	if existing, err := models.FindUserByEmail(app, DemoEmail); err == nil {
		return existing, nil
	}

	user, err := models.CreateUser(app)
	if err != nil {
		return nil, err
	}

	user.SetEmail(DemoEmail)
	user.SetPassword(DemoPassword)
	user.SetName("Demo User")
	user.SetVerified(true)

	if err := app.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

func ensureDemoPosts(app core.App, owner *models.User) error {
	demoPosts := []struct {
		title     string
		content   string
		published bool
	}{
		{"Welcome to the template", "<p>This post was created by the seeder.</p>", true},
		{"Working with services", "<p>Business logic lives in internal/services.</p>", true},
		{"An unpublished draft", "<p>Only the owner can see this one.</p>", false},
	}

	for _, p := range demoPosts {
		existing, _ := app.FindFirstRecordByData("posts", "title", p.title)
		if existing != nil {
			continue
		}

		post, err := models.CreatePost(app)
		if err != nil {
			return err
		}

		post.Apply(models.PostData{
			Title:     p.title,
			Content:   p.content,
			Published: p.published,
			OwnerID:   owner.Id,
		})

		// Save goes through the app so the slug hook fires too.
		if err := app.Save(post); err != nil {
			return err
		}
	}

	return nil
}
