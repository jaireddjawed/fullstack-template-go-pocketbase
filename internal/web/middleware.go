package web

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	gonertia "github.com/romsar/gonertia/v2"
)

// AuthCookie is the httpOnly cookie holding the PocketBase auth token for
// web (Inertia) sessions. API clients keep using the Authorization header.
const AuthCookie = "pb_auth"

// loadAuthFromCookie resolves the auth cookie into e.Auth, so web routes get
// the same auth semantics as PocketBase's own header-based middleware.
func loadAuthFromCookie(e *core.RequestEvent) error {
	if e.Auth == nil {
		if cookie, err := e.Request.Cookie(AuthCookie); err == nil && cookie.Value != "" {
			if record, err := e.App.FindAuthRecordByToken(cookie.Value, core.TokenTypeAuth); err == nil {
				e.Auth = record
			}
		}
	}
	return e.Next()
}

// requireGuest redirects authenticated users away from guest-only pages.
func requireGuest(e *core.RequestEvent) error {
	if e.Auth != nil {
		return e.Redirect(http.StatusSeeOther, "/")
	}
	return e.Next()
}

// requireWebAuth redirects unauthenticated users to the login page.
// (The API equivalent is apis.RequireAuth(), which returns 401 instead.)
func requireWebAuth(e *core.RequestEvent) error {
	if e.Auth == nil {
		return e.Redirect(http.StatusSeeOther, "/login")
	}
	return e.Next()
}

// inertiaMiddleware adapts gonertia's net/http middleware (Inertia protocol:
// asset versioning, redirect status fix-ups, ...) to PocketBase's router.
func inertiaMiddleware(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var nextErr error

		handler := i.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			e.Response = w
			e.Request = r
			nextErr = e.Next()
		}))

		handler.ServeHTTP(e.Response, e.Request)
		return nextErr
	}
}

// setAuthCookie writes the session cookie. Set Secure: true behind HTTPS.
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearAuthCookie expires the session cookie.
func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
