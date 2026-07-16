// Package web serves the Inertia.js frontend: it renders the root template,
// bridges gonertia to PocketBase's router, and implements cookie-based
// session auth on top of PocketBase auth records.
package web

import (
	"fmt"
	"net/http"
	"os"

	gonertia "github.com/romsar/gonertia/v2"
)

const (
	// RootTemplate is the Go html/template that hosts the Inertia app.
	RootTemplate = "frontend/root.html"

	// HotFile is written by the Vite dev server (see frontend/vite.config.ts).
	// Its presence switches gonertia into hot-reload mode.
	HotFile = "frontend/.hot"

	// BuildDir is where `npm run build` puts production assets.
	BuildDir = "frontend/dist"
)

// NewInertia builds the gonertia instance with Laravel-style Vite support:
// dev mode when frontend/.hot exists (assets served by the Vite dev server),
// production mode otherwise (hashed assets resolved via dist/manifest.json).
func NewInertia() (*gonertia.ViteInstance, error) {
	if !HasFrontendAssets() {
		return nil, fmt.Errorf("frontend assets unavailable: run `cd frontend && npm install`, then `make dev`, or build assets with `cd frontend && npm run build`")
	}

	base, err := gonertia.NewFromFile(RootTemplate)
	if err != nil {
		return nil, err
	}

	return gonertia.NewVite(base,
		gonertia.WithHotFile(HotFile),
		gonertia.WithBuildManifest(BuildDir+"/manifest.json"),
		gonertia.WithFallbackManifest(BuildDir+"/.vite/manifest.json"),
		gonertia.WithBuildDir("/build/"),
	)
}

// HasFrontendAssets reports whether gonertia can render asset tags without
// failing mid-response.
func HasFrontendAssets() bool {
	return hasHotFile() || HasBuild()
}

// HasBuild reports whether a production frontend build exists.
func HasBuild() bool {
	return fileExists(BuildDir+"/manifest.json") || fileExists(BuildDir+"/.vite/manifest.json")
}

func hasHotFile() bool {
	return fileExists(HotFile)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// renderError re-renders a component with Inertia validation errors
// (surfaces in the frontend as the `errors` page prop).
func renderError(i *gonertia.ViteInstance, w http.ResponseWriter, r *http.Request, component string, errs gonertia.ValidationErrors) error {
	ctx := gonertia.SetValidationErrors(r.Context(), errs)
	return i.Render(w, r.WithContext(ctx), component, nil)
}
