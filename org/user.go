package org

import (
	"context"
	"errors"
	"time"

	"github.com/cristianuser/go-bun-webserver/bunapp"
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

type Session struct {
	bun.BaseModel `bun:",alias:s"`

	ID             uint64                 `json:"-" bun:",pk,autoincrement"`
	UserId         uint64                 `json:"userId"`
	User           User                   `json:"user" bun:"rel:belongs-to"`
	Token          string                 `bun:",notnull,unique" json:"token,omitempty"`
	Provider       string                 `json:"provider" bun:"default:'LOCAL'"`
	LastTimeActive time.Time              `json:"lastTimeActive" bun:"default:current_timestamp"`
	ExpiresAt      time.Time              `json:"expiresAt"`
	DeviceInfo     map[string]interface{} `json:"deviceInfo" bun:"type:jsonb,default:'{}'"`
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

func (u *User) CreateSession(app *bunapp.App) (Session, error) {
	var session Session
	tokenTtl := 24 * time.Hour

	token, err := CreateUserToken(app, u.ID, tokenTtl)
	if err != nil {
		return session, err
	}
	session = Session{
		UserId:    u.ID,
		Token:     token,
		Provider:  "LOCAL",
		ExpiresAt: time.Now().Add(tokenTtl),
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

func (s *Session) UpdateLastTimeActive(ctx context.Context, app *bunapp.App) error {
	s.LastTimeActive = time.Now()
	if _, err := app.DB().NewUpdate().
		Model(s).
		Set("last_time_active = ?", s.LastTimeActive).
		Where("id = ?", s.ID).
		Exec(ctx); err != nil {
		return err
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

func SelectSessionByToken(ctx context.Context, app *bunapp.App, token string) (*Session, error) {
	session := new(Session)
	if err := app.DB().NewSelect().
		Model(session).
		Where("token = ?", token).
		Scan(ctx); err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}
	return session, nil
}
