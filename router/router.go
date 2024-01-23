package router

import (
	"github.com/devldm/go-server-rss/db"
	"github.com/devldm/go-server-rss/handlers"
	"github.com/devldm/go-server-rss/internal/database"
	"github.com/devldm/go-server-rss/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter(dbq *database.Queries) *chi.Mux {
	router := chi.NewRouter()
	config := db.NewAPIConfig(dbq)
	router.Use(middleware.ConfigMiddleware(config))

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlers.HandlerReadiness)

	v1Router.Get("/error", handlers.HandlerError)
	// users
	v1Router.Post("/users", handlers.HandlerCreateUser)

	v1Router.Group(func(v1Router chi.Router) {
		v1Router.Get("/users", middleware.MiddlewareAuth(handlers.HandlerGetUserByApiKey))
		// feeds
		v1Router.Post("/feeds", middleware.MiddlewareAuth(handlers.HandlerCreateFeed))
		// feed_follows
		v1Router.Post("/feed_follows", middleware.MiddlewareAuth(handlers.HandlerCreateFeedFollow))
		v1Router.Get("/feed_follows", middleware.MiddlewareAuth(handlers.HandlerGetFeedFollows))
		v1Router.Delete("/feed_follows/{feedFollowID}", middleware.MiddlewareAuth(handlers.HandlerDeleteFeedFollow))

		// user feed
		v1Router.Get("/posts", middleware.MiddlewareAuth(handlers.HandlerGetPostsForUser))

	})

	v1Router.Get("/feeds", handlers.HandlerGetFeeds)

	router.Mount("/v1", v1Router)

	return router
}
