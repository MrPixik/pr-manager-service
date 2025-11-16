package server

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"service-order-avito/internal/http/middleware"
	"service-order-avito/internal/http/server/handlers"
)

type TeamHandler interface {
	AddTeam(http.ResponseWriter, *http.Request)
	GetTeam(http.ResponseWriter, *http.Request)
	GetTeamStats(http.ResponseWriter, *http.Request)
}

type UserHandler interface {
	SetIsActive(http.ResponseWriter, *http.Request)
	GetReviewPullRequests(http.ResponseWriter, *http.Request)
}

type PullRequestHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Merge(http.ResponseWriter, *http.Request)
	ReassignReviewer(http.ResponseWriter, *http.Request)
}

func InitRouter(log *slog.Logger,
	teamHandler TeamHandler,
	userHandler UserHandler,
	prHandler PullRequestHandler,
) chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.WithLogger(log),
	)

	router.Get("/ping", handlers.PingGetHandler)
	router.Head("/healthcheck", handlers.HealthcheckHeadHandler)

	router.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.AddTeam)
		r.Get("/get", teamHandler.GetTeam)
		r.Get("/stats", teamHandler.GetTeamStats)
	})

	router.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetIsActive)
		r.Get("/getReview", userHandler.GetReviewPullRequests)
	})

	router.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", prHandler.Create)
		r.Post("/merge", prHandler.Merge)
		r.Post("/reassign", prHandler.ReassignReviewer)
	})
	return router
}
