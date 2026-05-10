package e2e

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const testsTimeout = 10 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	log := logger.FromContext(suiteCtx)
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	if err != nil {
		log.Error("Не удалось загрузить .env файл", slog.String("err", err.Error()))
		os.Exit(1)
	}

	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	log.Info("Запуск тестового окружения...")
	env, err = setupTestEnvironment(suiteCtx)
	if err != nil {
		log.Error("Не удалось поднять тестовое окружение", slog.String("err", err.Error()))
		os.Exit(1)
	}
})

var _ = ginkgo.AfterSuite(func() {
	log := logger.FromContext(context.Background())
	log.Info("Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
