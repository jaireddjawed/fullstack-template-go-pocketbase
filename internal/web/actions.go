package web

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	gonertia "github.com/romsar/gonertia/v2"

	"github.com/jaireddjawed/fullstack-template-golang/internal/models"
	"github.com/jaireddjawed/fullstack-template-golang/internal/services"
	"github.com/jaireddjawed/fullstack-template-golang/internal/types"
)

// Page actions render Inertia components. Each handler builds typed props
// (structs from internal/types, mirrored to TS by `make types`) and hands
// them to gonertia.

func sharedProps(e *core.RequestEvent) gonertia.Props {
	props := gonertia.Props{"auth": nil}

	if e.Auth != nil {
		user := models.NewUser(e.Auth)
		props["auth"] = types.AuthUser{
			ID:    user.Id,
			Email: user.Email(),
			Name:  user.Name(),
		}
	}

	return props
}

func home(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		stats, err := services.NewPostService(e.App).Stats()
		if err != nil {
			return err
		}

		props := sharedProps(e)
		props["stats"] = stats

		return i.Render(e.Response, e.Request, "Home", props)
	}
}

func postsIndex(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		records, err := e.App.FindRecordsByFilter(
			"posts",
			"published = true || owner = {:auth}",
			"-created",
			100,
			0,
			map[string]any{"auth": authID(e)},
		)
		if err != nil {
			return err
		}

		posts := make([]types.PostSummary, 0, len(records))
		for _, record := range records {
			post := models.NewPost(record)
			data := post.Data()
			posts = append(posts, types.PostSummary{
				ID:        post.Id,
				Title:     data.Title,
				Slug:      data.Slug,
				Content:   data.Content,
				Published: data.Published,
				IsOwner:   post.IsOwnedBy(authID(e)),
				Created:   post.GetDateTime("created").String(),
			})
		}

		props := sharedProps(e)
		props["posts"] = posts

		return i.Render(e.Response, e.Request, "Posts/Index", props)
	}
}

func loginPage(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return i.Render(e.Response, e.Request, "Auth/Login", sharedProps(e))
	}
}

func login(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		form := struct {
			Email    string `json:"email" form:"email"`
			Password string `json:"password" form:"password"`
		}{}
		if err := e.BindBody(&form); err != nil {
			return err
		}

		user, err := models.FindUserByEmail(e.App, form.Email)
		if err != nil || !user.ValidatePassword(form.Password) {
			return renderError(i, e.Response, e.Request, "Auth/Login", gonertia.ValidationErrors{
				"email": "Invalid email or password.",
			})
		}
		if !user.Verified() {
			return renderError(i, e.Response, e.Request, "Auth/Login", gonertia.ValidationErrors{
				"email": "Verify your email before logging in.",
			})
		}

		token, err := user.NewAuthToken()
		if err != nil {
			return err
		}

		setAuthCookie(e.Response, token)
		i.Redirect(e.Response, e.Request, "/")
		return nil
	}
}

func signupPage(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return i.Render(e.Response, e.Request, "Auth/Signup", sharedProps(e))
	}
}

func signup(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		form := struct{ Name, Email, Password, PasswordConfirm string }{}
		if err := e.BindBody(&form); err != nil {
			return err
		}
		if form.Email == "" || form.Password == "" || form.Password != form.PasswordConfirm {
			return renderError(i, e.Response, e.Request, "Auth/Signup", gonertia.ValidationErrors{"email": "Provide an email and matching passwords."})
		}
		user, err := models.CreateUser(e.App)
		if err != nil {
			return err
		}
		user.SetEmail(form.Email)
		user.SetPassword(form.Password)
		user.SetName(form.Name)
		if err := e.App.Save(user); err != nil {
			return renderError(i, e.Response, e.Request, "Auth/Signup", gonertia.ValidationErrors{"email": "We couldn't create that account. Check your details and try again."})
		}
		if err := mails.SendRecordVerification(e.App, user.Record); err != nil {
			e.App.Logger().Error("send verification email", "error", err)
		}
		i.Redirect(e.Response, e.Request, "/verify-email")
		return nil
	}
}

func verifyEmailPage(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return i.Render(e.Response, e.Request, "Auth/VerifyEmail", sharedProps(e))
	}
}

func requestVerification(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		form := struct{ Email string }{}
		if err := e.BindBody(&form); err != nil {
			return err
		}
		if user, err := models.FindUserByEmail(e.App, form.Email); err == nil && !user.Verified() {
			if err := mails.SendRecordVerification(e.App, user.Record); err != nil {
				e.App.Logger().Error("send verification email", "error", err)
			}
		}
		return i.Render(e.Response, e.Request, "Auth/VerifyEmail", gonertia.Props{"success": "If an account exists for that email, a verification link is on its way."})
	}
}

func confirmVerificationPage(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return i.Render(e.Response, e.Request, "Auth/ConfirmEmail", gonertia.Props{"token": e.Request.URL.Query().Get("token")})
	}
}

func confirmVerification(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		form := struct{ Token string }{}
		if err := e.BindBody(&form); err != nil {
			return err
		}
		user, err := e.App.FindAuthRecordByToken(form.Token, core.TokenTypeVerification)
		if err != nil {
			return renderError(i, e.Response, e.Request, "Auth/ConfirmEmail", gonertia.ValidationErrors{"token": "This verification link is invalid or has expired."})
		}
		user.SetVerified(true)
		if err := e.App.Save(user); err != nil {
			return err
		}
		i.Redirect(e.Response, e.Request, "/login")
		return nil
	}
}

func logout(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		clearAuthCookie(e.Response)
		i.Redirect(e.Response, e.Request, "/login")
		return nil
	}
}

func publishPost(i *gonertia.ViteInstance) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		_, err := services.NewPostService(e.App).Publish(e.Request.PathValue("id"), authID(e))
		if err != nil {
			return e.Redirect(http.StatusSeeOther, "/posts")
		}

		i.Back(e.Response, e.Request)
		return nil
	}
}

func authID(e *core.RequestEvent) string {
	if e.Auth == nil {
		return ""
	}
	return e.Auth.Id
}
