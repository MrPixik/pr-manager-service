package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service-order-avito/internal/config"
	"service-order-avito/internal/http/server"
	"service-order-avito/internal/http/server/handlers/pull_request"
	"service-order-avito/internal/http/server/handlers/team"
	"service-order-avito/internal/http/server/handlers/user"
	"service-order-avito/internal/repository/postgres"
	"service-order-avito/internal/service"
	"service-order-avito/pkg/logger"
	"syscall"
)

func main() {
	// CONFIG
	cfg := config.MustLoad()

	// LOGGER
	log := logger.MustInit(cfg.Env)
	log.Info("logger initialized")

	// GS context
	ctxApp, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// context для бд, отменяется после завершения сервера
	ctxDB, cancelDB := context.WithCancel(context.Background())
	defer cancelDB()

	// Conn to Postgres
	conn, err := postgres.ConnectPostgres(ctxDB, cfg.Postgres, cfg.Env)
	if err != nil {
		log.Error("init repository: " + err.Error())
		os.Exit(1)
	}
	defer func() {
		conn.Close()
		log.Info("connection with database closed")
	}()

	// Repository's Lay
	userRepo := postgres.NewUserRepositoryPostgres(conn)
	teamRepo := postgres.NewTeamRepositoryPostgres(conn, userRepo)
	prRepo := postgres.NewPullRequestRepositoryPostgres(conn, teamRepo, userRepo)
	log.Info("repository's lay initialized")

	// Service lay
	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	prService := service.NewPullRequestService(prRepo)
	log.Info("service's lay initialized")

	// Controller's lay
	teamHandler := team.NewTeamHandler(teamService)
	userHandler := user.NewUserHandler(userService)
	prHandler := pull_request.NewPullRequestHandler(prService)
	log.Info("courier handler initialized")

	// ROUTER & SERVER
	r := server.InitRouter(log, teamHandler, userHandler, prHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: r,
		BaseContext: func(net.Listener) context.Context {
			return ctxApp
		},
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.ShutdownTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server start up: %s", err.Error())
		}
	}()
	log.Info("listening on: " + cfg.HTTP.Port)

	gracefulShutdownServer(ctxApp, cfg.HTTP, log, srv, cancelDB)
}

func gracefulShutdownServer(ctxApp context.Context, cfg config.HTTPServer, log *slog.Logger, srv *http.Server, cancelDB context.CancelFunc) {
	<-ctxApp.Done()
	log.Info("shutdown signal received. starting graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown: %s", err.Error())
	} else {
		log.Info("server gracefully stopped")
	}

	cancelDB()
}
