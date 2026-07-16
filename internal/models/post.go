// Package models contains typed wrappers around PocketBase records —
// the equivalent of Laravel's Eloquent models. Each model embeds
// core.BaseRecordProxy, so it *is* the underlying record: pass it directly
// to app.Save() and all hooks/validations fire as usual.
//
// Models give you compile-time safety over field names in Go code, while
// the schema itself stays defined in migrations.
package models

import (
	"github.com/pocketbase/pocketbase/core"
)

// interface guard
var _ core.RecordProxy = (*Post)(nil)

// Post wraps a record of the "posts" collection.
type Post struct {
	core.BaseRecordProxy
}

// PostData is the typed representation of a posts record.
//
// The tags document the matching PocketBase field names. Use Data to read a
// record and Apply to copy a complete value back before saving the Post.
type PostData struct {
	ID        string `pb:"id"`
	Title     string `pb:"title"`
	Slug      string `pb:"slug"`
	Content   string `pb:"content"`
	Published bool   `pb:"published"`
	OwnerID   string `pb:"owner"`
}

// NewPost wraps an existing posts record.
func NewPost(record *core.Record) *Post {
	p := &Post{}
	p.SetProxyRecord(record)
	return p
}

// CreatePost returns a fresh, unsaved posts model.
func CreatePost(app core.App) (*Post, error) {
	collection, err := app.FindCollectionByNameOrId("posts")
	if err != nil {
		return nil, err
	}
	return NewPost(core.NewRecord(collection)), nil
}

// FindPostByID loads a post by id.
func FindPostByID(app core.App, id string) (*Post, error) {
	record, err := app.FindRecordById("posts", id)
	if err != nil {
		return nil, err
	}
	return NewPost(record), nil
}

// Data returns a typed snapshot of the record fields.
func (p *Post) Data() PostData {
	return PostData{
		ID:        p.Id,
		Title:     p.GetString("title"),
		Slug:      p.GetString("slug"),
		Content:   p.GetString("content"),
		Published: p.GetBool("published"),
		OwnerID:   p.GetString("owner"),
	}
}

// Apply copies the writable PostData fields to the record. The record ID is
// managed by PocketBase and is intentionally not written.
func (p *Post) Apply(data PostData) {
	p.Set("title", data.Title)
	p.Set("slug", data.Slug)
	p.Set("content", data.Content)
	p.Set("published", data.Published)
	p.Set("owner", data.OwnerID)
}

// IsOwnedBy reports whether the post belongs to the given user id.
func (p *Post) IsOwnedBy(userID string) bool {
	return userID != "" && p.Data().OwnerID == userID
}
