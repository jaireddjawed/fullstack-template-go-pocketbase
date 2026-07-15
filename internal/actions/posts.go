// Package actions contains the HTTP handlers for custom routes — the
// equivalent of Laravel's single-action controllers. Actions parse the
// request, call a service, and translate service results/errors into
// HTTP responses. No business logic lives here.
package actions

import (
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/jaireddjawed/fullstack-template-golang/internal/services"
	"github.com/jaireddjawed/fullstack-template-golang/internal/types"
)

const appVersion = "0.1.0"

// Health handles GET /api/app/health.
func Health(e *core.RequestEvent) error {
	return e.JSON(http.StatusOK, types.HealthResponse{
		Status:  "ok",
		Version: appVersion,
	})
}

// PostStats handles GET /api/app/posts/stats.
func PostStats(e *core.RequestEvent) error {
	stats, err := services.NewPostService(e.App).Stats()
	if err != nil {
		return apis.NewInternalServerError("failed to compute post stats", err)
	}
	return e.JSON(http.StatusOK, stats)
}

// PublishPost handles POST /api/app/posts/{id}/publish.
// The route is bound with apis.RequireAuth(), so e.Auth is always set.
func PublishPost(e *core.RequestEvent) error {
	post, err := services.NewPostService(e.App).Publish(e.Request.PathValue("id"), e.Auth.Id)

	switch {
	case errors.Is(err, services.ErrNotOwner):
		return apis.NewForbiddenError("you can only publish your own posts", nil)
	case err != nil:
		return apis.NewNotFoundError("post not found", err)
	}

	return e.JSON(http.StatusOK, types.PublishPostResponse{
		ID:        post.Id,
		Slug:      post.GetString("slug"),
		Published: post.GetBool("published"),
	})
}
