// Package services contains the business logic layer. Services are plain
// structs around a core.App — they know nothing about HTTP, so they can be
// called from actions (HTTP), hooks, CLI commands, or tests.
package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"github.com/jaireddjawed/fullstack-template-golang/internal/models"
	"github.com/jaireddjawed/fullstack-template-golang/internal/types"
)

var (
	// ErrNotOwner is returned when a user tries to act on a post they don't own.
	ErrNotOwner = errors.New("post does not belong to the authenticated user")
)

// PostService owns all business logic around the posts collection.
type PostService struct {
	app core.App
}

func NewPostService(app core.App) *PostService {
	return &PostService{app: app}
}

// Stats returns aggregate counts over the posts collection.
func (s *PostService) Stats() (*types.PostStats, error) {
	total, err := s.app.CountRecords("posts")
	if err != nil {
		return nil, err
	}

	published, err := s.app.CountRecords("posts", dbx.HashExp{"published": true})
	if err != nil {
		return nil, err
	}

	return &types.PostStats{
		Total:     total,
		Published: published,
		Drafts:    total - published,
	}, nil
}

// Publish marks a post as published after verifying ownership.
func (s *PostService) Publish(postID, userID string) (*models.Post, error) {
	post, err := models.FindPostByID(s.app, postID)
	if err != nil {
		return nil, err
	}

	if !post.IsOwnedBy(userID) {
		return nil, ErrNotOwner
	}

	data := post.Data()
	data.Published = true
	post.Apply(data)
	if err := s.app.Save(post); err != nil {
		return nil, err
	}

	return post, nil
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts a title into a url-safe slug ("Hello, World!" -> "hello-world").
func Slugify(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = nonSlugChars.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}
