package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/blazee5/url-shortener-rest-api/internal/config"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/redirect"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/url/delete"
	"github.com/blazee5/url-shortener-rest-api/internal/http-server/handlers/url/save"
	sl "github.com/blazee5/url-shortener-rest-api/internal/lib/logger/slog"
	"github.com/blazee5/url-shortener-rest-api/internal/storage/mongodb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	client, err := mongodb.Run(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	dao, err := mongodb.NewDAO(client.DB)
	if err != nil {
		log.Error("failed to init DAO", sl.Err(err))
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Post("/url", save.New(log, dao))
	r.Get("/{alias}", redirect.New(log, dao))
	r.Delete("/{alias}", delete.Delete(log, dao))

	log.Info(fmt.Sprintf("starting server on %s:%s", cfg.Host, cfg.Port))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("failed while stopping server", sl.Err(err))
	}

	if err := client.DB.Disconnect(context.Background()); err != nil {
		log.Error("failed while stopping database", sl.Err(err))
	}
}
