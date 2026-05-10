package e2e

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	"github.com/m4kson/rocket-factory/platform/pkg/testcontainers"
	"github.com/m4kson/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/m4kson/rocket-factory/platform/pkg/testcontainers/mongo"
	"github.com/m4kson/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/m4kson/rocket-factory/platform/pkg/testcontainers/path"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	inventoryAppName    = "inventory-service"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	grpcPortKey = "GRPC_PORT"

	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) (*TestEnvironment, error) {
	log := logger.FromContext(ctx)

	log.Info("Подготовка тестового окружения...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		log.Error("не удалось создать общую сеть", slog.String("err", err.Error()))
		return nil, err
	}
	log.Info("✅ Сеть успешно создана")

	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(*log),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		log.Error("не удалось запустить контейнер MongoDB", slog.String("err", err.Error()))
		return nil, err
	}
	log.Info("Контейнер MongoDB успешно запущен")

	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		testcontainers.MongoHostKey:     testcontainers.MongoContainerName,
		testcontainers.MongoPortKey:     testcontainers.MongoPort,
		testcontainers.MongoDatabaseKey: mongoDatabase,
		testcontainers.MongoUsernameKey: mongoUsername,
		testcontainers.MongoPasswordKey: mongoPassword,
		testcontainers.MongoAuthDBKey:   getEnvWithLogging(ctx, testcontainers.MongoAuthDBKey),
		testcontainers.LoggerLevelKey:   getEnvWithLogging(ctx, testcontainers.LoggerLevelKey),
		testcontainers.LoggerAsJsonKey:  getEnvWithLogging(ctx, testcontainers.LoggerAsJsonKey),
		testcontainers.GrpcPortKey:      grpcPort,
		"GRPC_HOST":                     "0.0.0.0",
		"MONGO_URI":                     "mongodb://admin:admin@inventory-mongo-test:27017/inventory?authSource=admin",
	}

	waitStrategy := wait.ForListeningPort(string(nat.Port(grpcPort + "/tcp"))).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(*log),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		log.Error("не удалось запустить контейнер приложения", slog.String("err", err.Error()))
		return nil, err
	}
	log.Info("Контейнер приложения успешно запущен")

	log.Info("Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}, nil
}

func getEnvWithLogging(ctx context.Context, key string) string {
	log := logger.FromContext(ctx)
	value := os.Getenv(key)
	if value == "" {
		log.Warn("Переменная окружения не установлена", slog.String("key", key))
	}

	return value
}
