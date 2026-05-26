package app

import (
	"context"
	"log/slog"

	"github.com/go-faster/errors"
	"github.com/m4kson/rocket-factory/assembly/internal/config"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

type App struct {
	diContainer *diContainer
	log         *slog.Logger
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "Shutdown signal received")
	case err := <-errCh:
		slog.ErrorContext(ctx, "Component crashed, shutting down", slog.String("error", err.Error()))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(a.log)
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

func (a *App) runConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "consumer initialized")

	err := a.diContainer.OrderPaidConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
