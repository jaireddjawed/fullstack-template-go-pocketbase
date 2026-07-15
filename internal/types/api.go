// Package types holds the request/response DTOs for the custom API routes.
//
// These structs are the single source of truth for API payload shapes:
// `make types` runs tygo, which converts every exported type in this
// package to TypeScript in shared/types.gen.ts. Never hand-edit the
// generated file — change the Go struct and regenerate.
package types

// HealthResponse is returned by GET /api/app/health.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// PostStats is returned by GET /api/app/posts/stats.
type PostStats struct {
	Total     int64 `json:"total"`
	Published int64 `json:"published"`
	Drafts    int64 `json:"drafts"`
}

// PublishPostResponse is returned by POST /api/app/posts/{id}/publish.
type PublishPostResponse struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Published bool   `json:"published"`
}
