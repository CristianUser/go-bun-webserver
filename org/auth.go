package org

import (
	"context"
	"net/http"
	"strings"

	"github.com/cristianuser/go-bun-webserver/bunapp"
	"github.com/cristianuser/go-bun-webserver/httputil"
	"github.com/uptrace/bunrouter"
)

type (
	userCtxKey       struct{}
	userErrCtxKey    struct{}
	sessionCtxKey    struct{}
	sessionErrCtxKey struct{}
)

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userCtxKey{}).(*User)
	return user
}

func SessionFromContext(ctx context.Context) *Session {
	session, _ := ctx.Value(sessionCtxKey{}).(*Session)
	return session
}

func authToken(req bunrouter.Request) string {
	const prefix = "Token "
	v := req.Header.Get("Authorization")
	v = strings.TrimPrefix(v, prefix)
	return v
}

type Middleware struct {
	app *bunapp.App
}

func NewMiddleware(app *bunapp.App) Middleware {
	return Middleware{
		app: app,
	}
}

func (m Middleware) User(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		ctx := req.Context()

		token := authToken(req)
		userID, err := decodeUserToken(m.app, token)
		if err != nil {
			ctx = context.WithValue(ctx, userErrCtxKey{}, err)
			return next(w, req.WithContext(ctx))
		}

		user, err := SelectUser(ctx, m.app, userID)
		if err != nil {
			ctx = context.WithValue(ctx, userErrCtxKey{}, err)
			return next(w, req.WithContext(ctx))
		}

		session, err := SelectSessionByToken(ctx, m.app, token)
		if err != nil {
			ctx = context.WithValue(ctx, sessionErrCtxKey{}, err)
			return next(w, req.WithContext(ctx))
		}

		err = session.UpdateLastTimeActive(ctx, m.app)
		if err != nil {
			ctx = context.WithValue(ctx, sessionErrCtxKey{}, err)
			return next(w, req.WithContext(ctx))
		}

		ctx = context.WithValue(ctx, userCtxKey{}, user)
		ctx = context.WithValue(ctx, sessionCtxKey{}, session)
		return next(w, req.WithContext(ctx))
	}
}

func (m Middleware) MustUser(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		if _, ok := req.Context().Value(userErrCtxKey{}).(error); ok {
			return httputil.JSON(w, bunrouter.H{
				"ok":      false,
				"message": "Invalid token",
			}, http.StatusUnauthorized)
		}
		if _, ok := req.Context().Value(sessionErrCtxKey{}).(error); ok {
			return httputil.JSON(w, bunrouter.H{
				"ok":      false,
				"message": "Invalid token",
			}, http.StatusUnauthorized)
		}

		return next(w, req)
	}
}
