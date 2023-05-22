package org

import (
	"context"
	"net/http"
	"time"

	"github.com/cristianuser/go-bun-webserver/bunapp"
	"github.com/cristianuser/go-bun-webserver/httputil"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	bun.BaseModel `bun:",alias:u"`

	ID       uint64 `json:"-" bun:",pk,autoincrement"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Username string `json:"username" bun:",unique"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Password string `bun:",notnull" json:"password,omitempty"`
}

type FollowUser struct {
	bun.BaseModel `bun:"alias:fu"`

	UserID         uint64
	FollowedUserID uint64
}

type Profile struct {
	bun.BaseModel `bun:"users,alias:u"`

	ID        uint64 `json:"-"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `bun:",scanonly" json:"following"`
}

func (u *User) CreateSession(app *bunapp.App, r *http.Request) (Session, error) {
	var session Session
	tokenTtl := 24 * time.Hour

	token, err := CreateUserToken(app, u.ID, tokenTtl)
	if err != nil {
		return session, err
	}

	ip := httputil.GetUserIP(r)
	browser, os := httputil.GetUserBrowserAndOS(r)
	session = Session{
		UserId:    u.ID,
		Token:     token,
		Provider:  "LOCAL",
		ExpiresAt: time.Now().Add(tokenTtl),
		DeviceInfo: map[string]interface{}{
			"ip":      ip,
			"browser": browser,
			"os":      os,
		},
	}

	if _, err := app.DB().NewInsert().
		Model(&session).
		Exec(app.Context()); err != nil {
		return session, err
	}
	return session, nil
}

func (u *User) ComparePassword(pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		return errUserNotFound
	}
	return nil
}

func NewProfile(user *User) *Profile {
	return &Profile{
		Username: user.Username,
		Image:    user.Image,
	}
}

func SelectUser(ctx context.Context, app *bunapp.App, id uint64) (*User, error) {
	user := new(User)
	if err := app.DB().NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func SelectUserByUsername(ctx context.Context, app *bunapp.App, username string) (*User, error) {
	user := new(User)
	if err := app.DB().NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx); err != nil {
		return nil, err
	}

	return user, nil
}
