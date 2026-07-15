package services_test

import (
	"errors"
	"testing"

	"github.com/jaireddjawed/fullstack-template-golang/internal/services"
	"github.com/jaireddjawed/fullstack-template-golang/internal/testutil"
)

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Hello, World!":        "hello-world",
		"  spaces  everywhere ": "spaces-everywhere",
		"Already-Slugged":      "already-slugged",
		"数字123 mixed":          "123-mixed",
	}

	for input, want := range cases {
		if got := services.Slugify(input); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestStats(t *testing.T) {
	app := testutil.NewApp(t)
	user := testutil.CreateUser(t, app, "stats@example.com")

	testutil.CreatePost(t, app, user.Id, "Published one", true)
	testutil.CreatePost(t, app, user.Id, "Published two", true)
	testutil.CreatePost(t, app, user.Id, "Draft", false)

	stats, err := services.NewPostService(app).Stats()
	if err != nil {
		t.Fatalf("Stats() returned error: %v", err)
	}

	if stats.Total != 3 || stats.Published != 2 || stats.Drafts != 1 {
		t.Errorf("Stats() = %+v, want total=3 published=2 drafts=1", stats)
	}
}

func TestPublishSetsPublished(t *testing.T) {
	app := testutil.NewApp(t)
	owner := testutil.CreateUser(t, app, "owner@example.com")
	post := testutil.CreatePost(t, app, owner.Id, "My Draft Post", false)

	updated, err := services.NewPostService(app).Publish(post.Id, owner.Id)
	if err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	if !updated.GetBool("published") {
		t.Error("Publish() did not set published=true")
	}
}

func TestPublishRejectsNonOwner(t *testing.T) {
	app := testutil.NewApp(t)
	owner := testutil.CreateUser(t, app, "owner@example.com")
	other := testutil.CreateUser(t, app, "other@example.com")
	post := testutil.CreatePost(t, app, owner.Id, "Owner's Post", false)

	_, err := services.NewPostService(app).Publish(post.Id, other.Id)
	if !errors.Is(err, services.ErrNotOwner) {
		t.Errorf("Publish() by non-owner returned %v, want ErrNotOwner", err)
	}
}

func TestSlugHookFiresOnCreate(t *testing.T) {
	app := testutil.NewApp(t)
	user := testutil.CreateUser(t, app, "hook@example.com")

	post := testutil.CreatePost(t, app, user.Id, "Hook Generated Slug!", false)

	if got := post.GetString("slug"); got != "hook-generated-slug" {
		t.Errorf("slug hook produced %q, want %q", got, "hook-generated-slug")
	}
}
