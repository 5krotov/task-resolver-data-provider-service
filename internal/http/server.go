package http

import (
	"context"
	"data-provider-service/internal/config"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Config config.HTTPConfig
	Server *http.Server
	Router *mux.Router
}

func NewServer(cfg config.HTTPConfig) *Server {
	return &Server{
		Config: cfg,
		Server: nil,
		Router: mux.NewRouter(),
	}
}

func (s *Server) Run() {
	log.Printf("running http server on %v ...\n", s.Config.Addr)
	s.Server = &http.Server{
		Addr:    s.Config.Addr,
		Handler: s.Router,
	}

	if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen and serve http server faild: %v", err)
	}
}

func (s *Server) Stop() {
	if s.Server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown http server failed: %v", err)
		}
	}
}
