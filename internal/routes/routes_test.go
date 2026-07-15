package routes_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tests"

	"github.com/jaireddjawed/fullstack-template-golang/internal/testutil"
)

// API tests use PocketBase's ApiScenario helper: it runs the request
// through the full router — including our custom routes, since
// testutil.NewApp binds them — and asserts on status code and body.
//
// DisableTestAppCleanup is set because testutil.NewApp already registers
// cleanup via t.Cleanup.

func TestHealthEndpoint(t *testing.T) {
	scenario := tests.ApiScenario{
		Name:                  "health returns ok",
		Method:                "GET",
		URL:                   "/api/app/health",
		ExpectedStatus:        200,
		ExpectedContent:       []string{`"status":"ok"`},
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return testutil.NewApp(tb)
		},
	}

	scenario.Test(t)
}

func TestPostStatsEndpoint(t *testing.T) {
	scenario := tests.ApiScenario{
		Name:                  "stats counts published and drafts",
		Method:                "GET",
		URL:                   "/api/app/posts/stats",
		ExpectedStatus:        200,
		ExpectedContent:       []string{`"total":2`, `"published":1`, `"drafts":1`},
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			app := testutil.NewApp(tb)
			user := testutil.CreateUser(tb, app, "stats@example.com")
			testutil.CreatePost(tb, app, user.Id, "Published", true)
			testutil.CreatePost(tb, app, user.Id, "Draft", false)
			return app
		},
	}

	scenario.Test(t)
}

func TestPublishRequiresAuth(t *testing.T) {
	scenario := tests.ApiScenario{
		Name:                  "publish without token is rejected",
		Method:                "POST",
		URL:                   "/api/app/posts/some-id/publish",
		ExpectedStatus:        401,
		ExpectedContent:       []string{`"data":{}`},
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return testutil.NewApp(tb)
		},
	}

	scenario.Test(t)
}

func TestPublishAsOwner(t *testing.T) {
	app := testutil.NewApp(t)
	owner := testutil.CreateUser(t, app, "owner@example.com")
	post := testutil.CreatePost(t, app, owner.Id, "My Draft", false)

	scenario := tests.ApiScenario{
		Name:            "owner can publish their draft",
		Method:          "POST",
		URL:             "/api/app/posts/" + post.Id + "/publish",
		ExpectedStatus:  200,
		ExpectedContent: []string{`"published":true`, `"slug":"my-draft"`},
		Headers: map[string]string{
			"Authorization": testutil.AuthToken(t, owner),
		},
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return app
		},
	}

	scenario.Test(t)
}

func TestPublishAsNonOwnerIsForbidden(t *testing.T) {
	app := testutil.NewApp(t)
	owner := testutil.CreateUser(t, app, "owner@example.com")
	other := testutil.CreateUser(t, app, "other@example.com")
	post := testutil.CreatePost(t, app, owner.Id, "Owner Post", false)

	scenario := tests.ApiScenario{
		Name:            "non-owner cannot publish",
		Method:          "POST",
		URL:             "/api/app/posts/" + post.Id + "/publish",
		ExpectedStatus:  403,
		ExpectedContent: []string{`You can only publish your own posts.`},
		Headers: map[string]string{
			"Authorization": testutil.AuthToken(t, other),
		},
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return app
		},
	}

	scenario.Test(t)
}
