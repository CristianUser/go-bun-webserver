package org

import (
	"context"
	"errors"
	"time"

	"github.com/cristianuser/go-bun-webserver/bunapp"
	"github.com/uptrace/bun"
)

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

func (s *Session) Destroy(ctx context.Context, app *bunapp.App) error {
	if _, err := app.DB().NewDelete().
		Model(s).
		Where("id = ?", s.ID).
		Exec(ctx); err != nil {
		return err
	}
	return nil
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
		session.Destroy(ctx, app)
		return nil, errors.New("Session Expired")
	}
	return session, nil
}
