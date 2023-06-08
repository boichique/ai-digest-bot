package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"digest_bot_database/internal/config"
	"digest_bot_database/internal/echox"
	"digest_bot_database/internal/log"
	"digest_bot_database/internal/modules/sources"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
)

const (
	dbConnectTimeout = 10 * time.Second
)

type Server struct {
	e       *echo.Echo
	cfg     *config.Config
	closers []func() error
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	logger, err := log.SetupLogger(cfg.Local, cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("setup logger: %w", err)
	}

	slog.SetDefault(logger)
	var closers []func() error
	db, err := getDB(context.Background(), cfg.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	closers = append(closers, func() error {
		db.Close()
		return nil
	})

	e := echo.New()
	sourcesModule := sources.NewModule(db)
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true

	api := e.Group("/api")
	api.Use(echox.Logger)

	api.PUT("/users/:userID", sourcesModule.Handler.CreateSource)
	api.GET("/users", sourcesModule.Handler.GetUsersIDList)
	api.GET("/users/:userID", func(c echo.Context) error {
		c.Set("YoutubeApiToken", cfg.YoutubeApiToken)
		return sourcesModule.Handler.GetNewVideosForUserSources(c)
	})
	api.GET("/users/:userID/digest", func(c echo.Context) error {
		c.Set("YoutubeApiToken", cfg.YoutubeApiToken)
		c.Set("ChatGPTApiToken", cfg.ChatGPTApiToken)
		return sourcesModule.Handler.GetDigestForUserSource(c)
	})
	api.GET("/users/:userID/sources", sourcesModule.Handler.GetUserSourcesByID)
	api.DELETE("/users/:userID", sourcesModule.Handler.DeleteSourceByLink)

	return &Server{
		e:       e,
		cfg:     cfg,
		closers: closers,
	}, nil
}

func (s *Server) Start() error {
	port := s.cfg.Port

	return s.e.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) Close() error {
	return withClosers(s.closers, nil)
}

func (s *Server) Port() (int, error) {
	listener := s.e.Listener
	if listener == nil {
		return 0, fmt.Errorf("server is not started")
	}

	addr := listener.Addr()
	if addr == nil {
		return 0, fmt.Errorf("server is not started")
	}

	return addr.(*net.TCPAddr).Port, nil
}

func getDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	return db, nil
}

func withClosers(closers []func() error, err error) error {
	errs := []error{err}

	for i := len(closers) - 1; i >= 0; i-- {
		if err := closers[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
