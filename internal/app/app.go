package app

import (
	"data-provider-service/internal/config"
	"data-provider-service/internal/http"
	"data-provider-service/internal/service/hello_world_service"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (*App) Run(cfg config.Config) {
	server := http.NewServer(cfg.HTTPConfig)
	service := hello_world_service.NewHelloWorldService()
	service.Register(server.Mux)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		server.Run()
	}()
	defer func() {
		server.Stop()
	}()

	<-stop
}
