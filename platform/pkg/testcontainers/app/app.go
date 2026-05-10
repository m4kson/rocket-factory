package app

import (
	"context"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/go-faster/errors"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultAppName        = "app"
	defaultAppPort        = "50051"
	defaultStartupTimeout = 1 * time.Minute
)

type Config struct {
	Name          string
	DockerfileDir string
	Dockerfile    string
	Port          string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Logger        *slog.Logger
}

type Container struct {
	container    testcontainers.Container
	externalHost string
	externalPort string
	cfg          *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := &Config{
		Name:          defaultAppName,
		Port:          defaultAppPort,
		Dockerfile:    "Dockerfile",
		DockerfileDir: ".",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultAppPort + "/tcp").WithStartupTimeout(defaultStartupTimeout),
		Env:           make(map[string]string),
		Logger:        logger.NewNop(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Name: cfg.Name,
		FromDockerfile: testcontainers.FromDockerfile{
			Context:        cfg.DockerfileDir,
			Dockerfile:     cfg.Dockerfile,
			BuildLogWriter: cfg.LogOutput,
		},
		Networks:           cfg.Networks,
		Env:                cfg.Env,
		WaitingFor:         cfg.StartupWait,
		ExposedPorts:       []string{cfg.Port + "/tcp"},
		HostConfigModifier: DefaultHostConfig(),
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start app genericContainer: %v", err)
	}

	mappedPort, err := genericContainer.MappedPort(ctx, cfg.Port+"/tcp")

	if err != nil {
		return nil, errors.Errorf("failed to get mapped externalPort: %v", err)
	}

	host, err := genericContainer.Host(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get genericContainer externalHost: %v", err)
	}

	go streamContainerLogs(ctx, genericContainer, cfg.LogOutput)

	cfg.Logger.Info("App container started", slog.String("uri:", net.JoinHostPort(host, mappedPort.Port())))

	return &Container{
		container:    genericContainer,
		externalHost: host,
		externalPort: mappedPort.Port(),
		cfg:          cfg,
	}, nil
}

func (a *Container) Address() string {
	return net.JoinHostPort(a.externalHost, a.externalPort)
}

func (a *Container) Terminate(ctx context.Context) error {
	return a.container.Terminate(ctx)
}

func streamContainerLogs(ctx context.Context, container testcontainers.Container, out io.Writer) {
	logs, err := container.Logs(ctx)
	if err != nil {
		slog.Error("failed to get container logs", slog.String("err", err.Error()))
		return
	}
	defer func() {
		err = logs.Close()
		if err != nil {
			slog.Error("failed to close container logs", slog.String("err", err.Error()))
		}
	}()

	go func() {
		_, err = io.Copy(out, logs)
		if err != nil && !errors.Is(err, io.EOF) {
			slog.Error("error copying container logs", slog.String("err", err.Error()))
		}
	}()
}

func DefaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = false
	}
}
