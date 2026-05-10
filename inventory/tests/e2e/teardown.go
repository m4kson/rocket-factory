package e2e

import (
	"context"
	"log/slog"

	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.FromContext(ctx)
	log.Info("Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	log.Info("Тестовое окружение успешно очищено")
}

func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.FromContext(ctx)
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			log.Error("не удалось остановить контейнер приложения", slog.String("err", err.Error()))
		} else {
			log.Info("Контейнер приложения остановлен")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			log.Error("не удалось остановить контейнер MongoDB", slog.String("err", err.Error()))
		} else {
			log.Info("Контейнер MongoDB остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			log.Error("не удалось удалить сеть", slog.String("err", err.Error()))
		} else {
			log.Info("Сеть удалена")
		}
	}
}
