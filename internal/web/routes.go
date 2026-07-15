package web

import (
	"os"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// Register mounts the Inertia web app on the PocketBase router.
//
// Returns an error if the root template is missing (e.g. in backend-only
// test runs) — the caller logs and continues, keeping /api/* fully usable.
func Register(se *core.ServeEvent) error {
	inertia, err := NewInertia()
	if err != nil {
		return err
	}

	// Production assets (built by `npm run build` into frontend/dist,
	// referenced as /build/... by the manifest resolver).
	if HasBuild() {
		se.Router.GET("/build/{path...}", apis.Static(os.DirFS(BuildDir), false))
	}

	pages := se.Router.Group("")
	pages.BindFunc(loadAuthFromCookie)
	pages.BindFunc(inertiaMiddleware(inertia))

	pages.GET("/", home(inertia))
	pages.GET("/posts", postsIndex(inertia)).BindFunc(requireWebAuth)
	pages.POST("/posts/{id}/publish", publishPost(inertia)).BindFunc(requireWebAuth)

	pages.GET("/login", loginPage(inertia)).BindFunc(requireGuest)
	pages.POST("/login", login(inertia)).BindFunc(requireGuest)
	pages.POST("/logout", logout(inertia))

	return nil
}
