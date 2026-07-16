// Package testutil provides helpers for backend tests: a fully migrated
// in-memory-ish test app (fresh temp data dir per test) with all hooks and
// routes bound, plus record factories.
package testutil

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"github.com/jaireddjawed/fullstack-template-go-pocketbase/internal/app"

	// Register app migrations so the test app is created fully migrated.
	_ "github.com/jaireddjawed/fullstack-template-go-pocketbase/migrations"
)

// NewApp returns a test app with all migrations applied and the project's
// hooks/routes bound — the same wiring production runs. Cleanup is
// registered automatically.
func NewApp(t testing.TB) *tests.TestApp {
	t.Helper()

	testApp, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create test app: %v", err)
	}
	t.Cleanup(testApp.Cleanup)

	app.Bind(testApp)

	return testApp
}

// CreateUser inserts a users record and returns it.
func CreateUser(t testing.TB, testApp *tests.TestApp, email string) *core.Record {
	t.Helper()

	users, err := testApp.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatalf("users collection missing: %v", err)
	}

	user := core.NewRecord(users)
	user.SetEmail(email)
	user.SetPassword("test-password-123")
	user.SetVerified(true)

	if err := testApp.Save(user); err != nil {
		t.Fatalf("failed to create user %s: %v", email, err)
	}
	return user
}

// CreatePost inserts a posts record owned by the given user and returns it.
func CreatePost(t testing.TB, testApp *tests.TestApp, ownerID, title string, published bool) *core.Record {
	t.Helper()

	posts, err := testApp.FindCollectionByNameOrId("posts")
	if err != nil {
		t.Fatalf("posts collection missing: %v", err)
	}

	post := core.NewRecord(posts)
	post.Set("title", title)
	post.Set("published", published)
	post.Set("owner", ownerID)

	if err := testApp.Save(post); err != nil {
		t.Fatalf("failed to create post %q: %v", title, err)
	}
	return post
}

// AuthToken returns a valid auth token for the record, for use in
// Authorization headers of API scenarios.
func AuthToken(t testing.TB, record *core.Record) string {
	t.Helper()

	token, err := record.NewAuthToken()
	if err != nil {
		t.Fatalf("failed to generate auth token: %v", err)
	}
	return token
}
