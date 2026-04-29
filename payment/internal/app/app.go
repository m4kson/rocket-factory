package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/m4kson/rocket-factory/payment/internal/config"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	"github.com/m4kson/rocket-factory/platform/pkg/grpc/health"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (a *App) Run() error {
	go func() {
		a.log.Info("starting grpc server", slog.String("address", a.listener.Addr().String()))
		if err := a.grpcServer.Serve(a.listener); err != nil {
			a.log.Error("grpc server error", slog.String("err", err.Error()))
		}
	}()

	<-closer.Done()
	return nil
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
		ServiceName: "payment",
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
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("GRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API())

	return nil
}
