package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m4kson/rocket-factory/order/internal/config"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	appMiddleware "github.com/m4kson/rocket-factory/platform/pkg/middleware"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/pkg/errors"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	log         *slog.Logger
}

func New(ctx context.Context) (*App, error) {
	app := &App{}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	return app, nil

}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	go func() {
		if err := a.startServer(); err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		a.log.InfoContext(ctx, "Shutdown signal received")
	case err := <-errCh:
		a.log.ErrorContext(ctx, "Component crashed, shutting down", slog.String("error", err.Error()))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDi,
		a.initMigrations,
		a.initCloser,
		a.initServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	log := logger.New(logger.Config{
		Level:       config.AppConfig().Logger.Level(),
		AsJson:      config.AppConfig().Logger.AsJson(),
		ServiceName: "order",
		Environment: "local", //todo add this ot config
		AddSource:   true,    //todo getEnv("ENV", "production") == "local"
	})

	a.log = log

	log.Info("logger initialized")

	return nil
}

func (a *App) initDi(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(a.log)
	return nil
}

func (a *App) initServer(ctx context.Context) error {
	handler, err := a.diContainer.OrderV1API(ctx)
	if err != nil {
		return err
	}

	ogenServer, err := orderV1.NewServer(handler)
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Use(appMiddleware.Logger(a.log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/api/v1", ogenServer)

	a.httpServer = &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", config.AppConfig().HttpServer.Port()),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       config.AppConfig().HttpServer.ReadHeaderTimeout(),
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	closer.AddNamed("http server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	a.log.Info("http server initialized",
		slog.String("addr", a.httpServer.Addr),
	)

	return nil
}

func (a *App) startServer() error {
	a.log.Info("server starting")
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("failed to start http server", slog.String("error", err.Error()), slog.String("port", config.AppConfig().HttpServer.Port()))
		return err
	}

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	runner, err := a.diContainer.MigrationRunner(ctx)
	if err != nil {
		a.log.Error("failed to create migration runner", slog.String("error", err.Error()))
		return err
	}
	err = runner.Up()
	if err != nil {
		a.log.Error("failed to run migrations", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	svc, err := a.diContainer.OrderConsumerService(ctx)
	if err != nil {
		return fmt.Errorf("runConsumer: %w", err)
	}

	a.log.InfoContext(ctx, "kafka consumer starting")
	return svc.RunConsumer(ctx)
}
