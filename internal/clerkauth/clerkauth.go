// Package clerkauth authenticates API requests with Clerk session tokens.
//
// The frontend sends Clerk JWTs in the Authorization header. The middleware
// verifies them (JWKS signature check via the Clerk SDK), maps the Clerk
// identity to a record in the "users" collection — provisioning one on
// first sight — and sets e.Auth. From there, everything downstream
// (collection API rules, apis.RequireAuth(), services) works exactly as
// with native PocketBase auth.
package clerkauth

import (
	"context"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"

	"github.com/jaireddjawed/fullstack-template-golang/internal/models"
)

// Identity is the verified Clerk identity extracted from a session token.
type Identity struct {
	ClerkUserID string
	Email       string
	Name        string
}

// Verifier validates a token and resolves it to a Clerk identity.
// It is an interface so tests can substitute a fake (see middleware_test.go).
type Verifier interface {
	Verify(ctx context.Context, token string) (*Identity, error)
}

// NewVerifierFromEnv returns the production verifier, or nil when
// CLERK_SECRET_KEY is not configured (Clerk auth disabled).
func NewVerifierFromEnv() Verifier {
	key := os.Getenv("CLERK_SECRET_KEY")
	if key == "" {
		return nil
	}
	clerk.SetKey(key)
	return &clerkVerifier{}
}

type clerkVerifier struct{}

// Verify checks the JWT signature against Clerk's JWKS and fetches the
// user's email/name from the Clerk management API (only reached on cache
// misses in practice, since provisioning happens once per user).
func (v *clerkVerifier) Verify(ctx context.Context, token string) (*Identity, error) {
	claims, err := clerkjwt.Verify(ctx, &clerkjwt.VerifyParams{Token: token})
	if err != nil {
		return nil, err
	}

	identity := &Identity{ClerkUserID: claims.RegisteredClaims.Subject}

	if u, err := clerkuser.Get(ctx, identity.ClerkUserID); err == nil {
		for _, email := range u.EmailAddresses {
			if u.PrimaryEmailAddressID != nil && email.ID == *u.PrimaryEmailAddressID {
				identity.Email = email.EmailAddress
			}
		}
		if u.FirstName != nil {
			identity.Name = *u.FirstName
		}
		if u.LastName != nil && identity.Name != "" {
			identity.Name += " " + *u.LastName
		}
	}

	return identity, nil
}

// Middleware resolves Clerk bearer tokens into e.Auth. Requests that are
// already authenticated (native PocketBase tokens) or carry no/invalid
// tokens pass through untouched — route guards decide what's required.
func Middleware(v Verifier) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Auth != nil {
			return e.Next()
		}

		token := e.Request.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		if token == "" {
			return e.Next()
		}

		identity, err := v.Verify(e.Request.Context(), token)
		if err != nil {
			// Not a (valid) Clerk token; continue unauthenticated.
			return e.Next()
		}

		user, err := ProvisionUser(e.App, identity)
		if err != nil {
			return err
		}

		e.Auth = user.ProxyRecord()
		return e.Next()
	}
}

// ProvisionUser returns the users record linked to the Clerk identity,
// creating (or linking by email) one on first sight.
func ProvisionUser(app core.App, identity *Identity) (*models.User, error) {
	if user, err := models.FindUserByClerkID(app, identity.ClerkUserID); err == nil {
		return user, nil
	}

	// Link an existing record with the same email (e.g. pre-Clerk data).
	if identity.Email != "" {
		if user, err := models.FindUserByEmail(app, identity.Email); err == nil {
			user.SetClerkID(identity.ClerkUserID)
			if err := app.Save(user); err != nil {
				return nil, err
			}
			return user, nil
		}
	}

	user, err := models.CreateUser(app)
	if err != nil {
		return nil, err
	}

	user.SetClerkID(identity.ClerkUserID)
	user.SetEmail(identity.Email)
	user.SetName(identity.Name)
	user.SetVerified(true)
	// Password auth is never used for Clerk-managed users; satisfy the
	// auth collection's requirement with a random secret.
	user.SetPassword(security.RandomString(32))

	if err := app.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}
