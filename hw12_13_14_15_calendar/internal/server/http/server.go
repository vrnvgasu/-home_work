package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	server http.Server
	logger Logger
	app    App
}

type Logger interface {
	Warn(msg string)
	Info(msg string)
	Error(msg string)
	File(msg string)
}

type App interface {
	CreateEvent(ctx context.Context, e storage.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, e storage.Event) error
	EventList(ctx context.Context, params storage.Params) ([]storage.Event, error)
}

func NewServer(logger Logger, app App) *Server {
	return &Server{
		server: http.Server{
			Addr:              fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port),
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	http.Handle("/hello", s.loggingMiddleware(HelloHandler{Logger: s.logger}))

	http.HandleFunc("POST /api/events", s.Add())
	http.HandleFunc("PUT /api/events", s.Update())
	http.HandleFunc("GET /api/events", s.List())

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
