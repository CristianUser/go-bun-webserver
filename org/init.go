package org

import (
	"context"

	"github.com/cristianuser/go-bun-webserver/bunapp"
)

func init() {
	bunapp.OnStart("org.initRoutes", func(ctx context.Context, app *bunapp.App) error {
		middleware := NewMiddleware(app)
		userHandler := NewUserHandler(app)
		authHandler := NewSessionHandler(app)

		app.Router().GET("/health", HealthCheckHandler)

		g := app.APIRouter().WithMiddleware(middleware.User)

		g.POST("/auth/login", authHandler.Login)
		g.POST("/auth/logout", authHandler.Logout)
		g.POST("/auth/register", userHandler.Create)

		g.POST("/users", userHandler.Create)
		g.GET("/profiles/:username", userHandler.Profile)

		g = g.WithMiddleware(middleware.MustUser)

		g.GET("/user/", userHandler.Current)
		g.PUT("/user/", userHandler.Update)

		g.POST("/profiles/:username/follow", userHandler.Follow)
		g.DELETE("/profiles/:username/follow", userHandler.Unfollow)

		return nil
	})
}
