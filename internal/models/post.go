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

func (p *Post) Title() string          { return p.GetString("title") }
func (p *Post) SetTitle(title string)  { p.Set("title", title) }
func (p *Post) Slug() string           { return p.GetString("slug") }
func (p *Post) SetSlug(slug string)    { p.Set("slug", slug) }
func (p *Post) Content() string        { return p.GetString("content") }
func (p *Post) SetContent(html string) { p.Set("content", html) }
func (p *Post) Published() bool        { return p.GetBool("published") }
func (p *Post) SetPublished(v bool)    { p.Set("published", v) }
func (p *Post) OwnerID() string        { return p.GetString("owner") }
func (p *Post) SetOwnerID(id string)   { p.Set("owner", id) }

// IsOwnedBy reports whether the post belongs to the given user id.
func (p *Post) IsOwnedBy(userID string) bool {
	return userID != "" && p.OwnerID() == userID
}
