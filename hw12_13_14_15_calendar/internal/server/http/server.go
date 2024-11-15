package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	server http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Warn(msg string)
	Info(msg string)
	Error(msg string)
	File(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
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
