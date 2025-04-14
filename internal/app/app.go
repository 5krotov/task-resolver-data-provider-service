package app

import (
	"context"
	"data-provider-service/internal/cache"
	"data-provider-service/internal/cache/repository/redis"
	"data-provider-service/internal/config"
	"data-provider-service/internal/handler/data_provider_handler"
	"data-provider-service/internal/http"
	"data-provider-service/internal/provider"
	"data-provider-service/internal/provider/repository/postgres"
	"data-provider-service/internal/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (*App) Run(cfg config.Config) error {
	server := http.NewServer(cfg.HTTPConfig)

	redisRepository, err := redis.NewRepository(cfg.RedisConfig)
	if err != nil {
		return fmt.Errorf("failed to create redis: %v", err)
	}
	dataCache := cache.NewCache(redisRepository)

	postgresRepository, err := postgres.NewRepository(cfg.PostgresConfig)
	if err != nil {
		return fmt.Errorf("failed to create postgres: %v", err)
	}
	dataProvider := provider.NewProvider(postgresRepository)

	dataProviderService := service.NewDataProviderService(dataCache, dataProvider)

	serviceHandler := data_provider_handler.NewDataProviderHandler(dataProviderService)

	serviceHandler.Register(server.Router)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server.Run()
	}()
	defer func() {
		server.Stop()
		redisRepository.Close(context.Background())
		postgresRepository.Close()
	}()

	<-stop

	return nil
}
