package types

// Props shared with Inertia pages. Regenerate the TS with `make types`.

// AuthUser is the authenticated user shared with every Inertia page.
type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// PostSummary is a post as rendered by Inertia pages.
type PostSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
	IsOwner   bool   `json:"isOwner"`
	Created   string `json:"created"`
}
