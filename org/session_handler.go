package org

import (
	"net/http"

	"github.com/cristianuser/go-bun-webserver/bunapp"
	"github.com/cristianuser/go-bun-webserver/httputil"
	"github.com/uptrace/bunrouter"
)

type SessionHandler struct {
	app *bunapp.App
}

func NewSessionHandler(app *bunapp.App) SessionHandler {
	return SessionHandler{
		app: app,
	}
}

func (h SessionHandler) Login(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := httputil.UnmarshalJSON(w, req, &credentials, 10<<kb); err != nil {
		return err
	}

	user := new(User)
	if err := h.app.DB().NewSelect().
		Model(user).
		Where("username = ? OR email = ?", credentials.Username, credentials.Username).
		Scan(ctx); err != nil {
		return err
	}

	if err := user.ComparePassword(credentials.Password); err != nil {
		return err
	}

	session, err := user.CreateSession(h.app, req.Request)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"user":  user,
		"token": session.Token,
	})
}

func (h SessionHandler) Logout(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	session := SessionFromContext(ctx)
	if session == nil {
		return nil
	}

	err := session.Destroy(ctx, h.app)
	if err != nil {
		return err
	}

	return nil
}
