package clerkauth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"github.com/jaireddjawed/fullstack-template-golang/internal/clerkauth"
	"github.com/jaireddjawed/fullstack-template-golang/internal/testutil"
)

// fakeVerifier accepts a single known token and returns a fixed identity.
type fakeVerifier struct {
	identity clerkauth.Identity
	token    string
}

func (f *fakeVerifier) Verify(_ context.Context, token string) (*clerkauth.Identity, error) {
	if token != f.token {
		return nil, errors.New("invalid token")
	}
	id := f.identity
	return &id, nil
}

func newFake() *fakeVerifier {
	return &fakeVerifier{
		token: "valid-clerk-token",
		identity: clerkauth.Identity{
			ClerkUserID: "user_clerk123",
			Email:       "clerk@example.com",
			Name:        "Clerk User",
		},
	}
}

// bindMiddleware attaches the clerk middleware inside an ApiScenario,
// mirroring what app.Bind does in production with the real verifier.
func bindMiddleware(v clerkauth.Verifier) func(testing.TB, *tests.TestApp, *core.ServeEvent) {
	return func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
		e.Router.BindFunc(clerkauth.Middleware(v))
	}
}

func TestProvisionCreatesUser(t *testing.T) {
	app := testutil.NewApp(t)
	fake := newFake()

	user, err := clerkauth.ProvisionUser(app, &fake.identity)
	if err != nil {
		t.Fatalf("ProvisionUser() returned error: %v", err)
	}

	if user.ClerkID() != "user_clerk123" || user.Email() != "clerk@example.com" {
		t.Errorf("provisioned user has wrong identity: clerk_id=%q email=%q", user.ClerkID(), user.Email())
	}
	if !user.Verified() {
		t.Error("provisioned user should be marked verified")
	}
}

func TestProvisionIsIdempotent(t *testing.T) {
	app := testutil.NewApp(t)
	fake := newFake()

	first, err := clerkauth.ProvisionUser(app, &fake.identity)
	if err != nil {
		t.Fatal(err)
	}
	second, err := clerkauth.ProvisionUser(app, &fake.identity)
	if err != nil {
		t.Fatal(err)
	}

	if first.Id != second.Id {
		t.Errorf("provisioning twice created two users: %s and %s", first.Id, second.Id)
	}
}

func TestProvisionLinksExistingUserByEmail(t *testing.T) {
	app := testutil.NewApp(t)
	existing := testutil.CreateUser(t, app, "clerk@example.com")
	fake := newFake()

	user, err := clerkauth.ProvisionUser(app, &fake.identity)
	if err != nil {
		t.Fatal(err)
	}

	if user.Id != existing.Id {
		t.Errorf("expected existing user %s to be linked, got new user %s", existing.Id, user.Id)
	}
	if user.ClerkID() != "user_clerk123" {
		t.Errorf("existing user was not linked to the clerk id, got %q", user.ClerkID())
	}
}

func TestMiddlewareAuthenticatesRequest(t *testing.T) {
	fake := newFake()

	scenario := tests.ApiScenario{
		Name:                  "clerk token resolves to a provisioned user",
		Method:                "GET",
		URL:                   "/api/app/me",
		Headers:               map[string]string{"Authorization": "Bearer valid-clerk-token"},
		ExpectedStatus:        200,
		ExpectedContent:       []string{`"email":"clerk@example.com"`, `"clerkId":"user_clerk123"`},
		BeforeTestFunc:        bindMiddleware(fake),
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return testutil.NewApp(tb)
		},
	}

	scenario.Test(t)
}

func TestMiddlewareRejectsInvalidToken(t *testing.T) {
	fake := newFake()

	scenario := tests.ApiScenario{
		Name:                  "invalid token stays unauthenticated",
		Method:                "GET",
		URL:                   "/api/app/me",
		Headers:               map[string]string{"Authorization": "Bearer wrong-token"},
		ExpectedStatus:        401,
		ExpectedContent:       []string{`"data":{}`},
		BeforeTestFunc:        bindMiddleware(fake),
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return testutil.NewApp(tb)
		},
	}

	scenario.Test(t)
}

func TestMiddlewareKeepsNativeAuth(t *testing.T) {
	app := testutil.NewApp(t)
	native := testutil.CreateUser(t, app, "native@example.com")
	fake := newFake()

	scenario := tests.ApiScenario{
		Name:                  "native pocketbase tokens still work",
		Method:                "GET",
		URL:                   "/api/app/me",
		Headers:               map[string]string{"Authorization": testutil.AuthToken(t, native)},
		ExpectedStatus:        200,
		ExpectedContent:       []string{`"email":"native@example.com"`},
		BeforeTestFunc:        bindMiddleware(fake),
		DisableTestAppCleanup: true,
		TestAppFactory: func(tb testing.TB) *tests.TestApp {
			return app
		},
	}

	scenario.Test(t)
}
