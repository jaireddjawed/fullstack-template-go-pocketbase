package actions

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"

	"github.com/jaireddjawed/fullstack-template-golang/internal/models"
	"github.com/jaireddjawed/fullstack-template-golang/internal/types"
)

// Me handles GET /api/app/me. Bound with apis.RequireAuth().
func Me(e *core.RequestEvent) error {
	user := models.NewUser(e.Auth)

	return e.JSON(http.StatusOK, types.AuthUser{
		ID:      user.Id,
		Email:   user.Email(),
		Name:    user.Name(),
		ClerkID: user.ClerkID(),
	})
}
