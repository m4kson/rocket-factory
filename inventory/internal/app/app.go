package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/m4kson/rocket-factory/inventory/internal/config"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	"github.com/m4kson/rocket-factory/platform/pkg/grpc/health"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDI,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
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
		ServiceName: "inventory",
		Environment: "local", //todo add this ot config
		AddSource:   true,    //todo getEnv("ENV", "production") == "local"
	})

	a.log = log

	log.Info("logger initialized")

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

func (a *App) initListener(_ context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig().Grpc.Port()))
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := lis.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = lis

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	closer.AddNamed("GRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	inventoryV1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	a.log.Info("starting gRPC server", "address", config.AppConfig().Grpc.Port())

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
