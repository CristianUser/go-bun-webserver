package org

import (
	"errors"
	"net/http"

	"github.com/uptrace/bunrouter"
	"golang.org/x/crypto/bcrypt"

	"github.com/cristianuser/go-bun-webserver/bunapp"
	"github.com/cristianuser/go-bun-webserver/httputil"
	"github.com/uptrace/bun"
)

const kb = 10

var errUserNotFound = errors.New("Not registered email or invalid password")

type UserHandler struct {
	app *bunapp.App
}

func NewUserHandler(app *bunapp.App) UserHandler {
	return UserHandler{
		app: app,
	}
}

func (*UserHandler) Current(w http.ResponseWriter, req bunrouter.Request) error {
	user := UserFromContext(req.Context())
	session := SessionFromContext(req.Context())
	return bunrouter.JSON(w, bunrouter.H{
		"user":    user,
		"session": session,
	})
}

func (h UserHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	var in struct {
		User *User `json:"user"`
	}

	if err := httputil.UnmarshalJSON(w, req, &in, 10<<kb); err != nil {
		return err
	}

	if in.User == nil {
		return errors.New(`JSON field "user" is required`)
	}

	user := in.User

	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	if _, err := h.app.DB().NewInsert().
		Model(user).
		Exec(ctx); err != nil {
		return err
	}

	user.Password = ""
	return bunrouter.JSON(w, bunrouter.H{
		"user": user,
	})
}

func (h UserHandler) Update(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	authUser := UserFromContext(ctx)

	var in struct {
		User *User `json:"user"`
	}

	if err := httputil.UnmarshalJSON(w, req, &in, 10<<kb); err != nil {
		return err
	}

	if in.User == nil {
		return errors.New(`JSON field "user" is required`)
	}

	user := in.User

	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	if _, err = h.app.DB().NewUpdate().
		Model(authUser).
		Set("email = ?", user.Email).
		Set("username = ?", user.Username).
		Set("password = ?", user.Password).
		Set("image = ?", user.Image).
		Where("id = ?", authUser.ID).
		Returning("*").
		Exec(ctx); err != nil {
		return err
	}

	user.Password = ""
	return bunrouter.JSON(w, bunrouter.H{
		"user": authUser,
	})
}

func (h UserHandler) Profile(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	followingColumn := func(q *bun.SelectQuery) *bun.SelectQuery {
		if authUser, ok := ctx.Value(userCtxKey{}).(*User); ok {
			subq := h.app.DB().NewSelect().
				Model((*FollowUser)(nil)).
				Where("fu.followed_user_id = u.id").
				Where("fu.user_id = ?", authUser.ID)

			q = q.ColumnExpr("EXISTS (?) AS following", subq)
		} else {
			q = q.ColumnExpr("false AS following")
		}

		return q
	}

	user := new(User)
	if err := h.app.DB().NewSelect().
		Model(user).
		ColumnExpr("u.*").
		Apply(followingColumn).
		Where("username = ?", req.Param("username")).
		Scan(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"profile": NewProfile(user),
	})
}

func (h UserHandler) Follow(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	authUser := UserFromContext(ctx)

	user, err := SelectUserByUsername(ctx, h.app, req.Param("username"))
	if err != nil {
		return err
	}

	followUser := &FollowUser{
		UserID:         authUser.ID,
		FollowedUserID: user.ID,
	}
	if _, err := h.app.DB().NewInsert().
		Model(followUser).
		Exec(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"profile": NewProfile(user),
	})
}

func (h UserHandler) Unfollow(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	authUser := UserFromContext(ctx)

	user, err := SelectUserByUsername(ctx, h.app, req.Param("username"))
	if err != nil {
		return err
	}

	if _, err := h.app.DB().NewDelete().
		Model((*FollowUser)(nil)).
		Where("user_id = ?", authUser.ID).
		Where("followed_user_id = ?", user.ID).
		Exec(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"profile": NewProfile(user),
	})
}

func hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
