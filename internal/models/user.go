package models

import (
	"github.com/pocketbase/pocketbase/core"
)

// interface guard
var _ core.RecordProxy = (*User)(nil)

// User wraps a record of the built-in "users" auth collection.
type User struct {
	core.BaseRecordProxy
}

// NewUser wraps an existing users record.
func NewUser(record *core.Record) *User {
	u := &User{}
	u.SetProxyRecord(record)
	return u
}

// CreateUser returns a fresh, unsaved users model.
func CreateUser(app core.App) (*User, error) {
	collection, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		return nil, err
	}
	return NewUser(core.NewRecord(collection)), nil
}

// FindUserByEmail loads a user by email.
func FindUserByEmail(app core.App, email string) (*User, error) {
	record, err := app.FindAuthRecordByEmail("users", email)
	if err != nil {
		return nil, err
	}
	return NewUser(record), nil
}

func (u *User) Name() string        { return u.GetString("name") }
func (u *User) SetName(name string) { u.Set("name", name) }
