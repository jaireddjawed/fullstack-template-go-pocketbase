package types

// AuthUser is returned by GET /api/app/me — the PocketBase user record the
// authenticated request resolved to (natively or via Clerk).
type AuthUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	ClerkID string `json:"clerkId"`
}
