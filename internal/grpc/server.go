package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"data-provider-service/internal/config"
	"data-provider-service/internal/service"
	"fmt"
	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type Server struct {
	config config.GRPCConfig
	agent  *DataProviderServiceAgent
	server *grpc.Server
}

func NewServer(cfg config.GRPCConfig, service *service.DataProviderService) (*Server, error) {
	var server *grpc.Server
	if cfg.UseTLS {
		serverCert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
		if err != nil {
			log.Fatalf("Failed to load server cert: %v", err)
		}

		caCert, err := os.ReadFile(cfg.CA)
		if err != nil {
			log.Fatalf("Failed to read CA cert: %v", err)
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)

		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to load server creds: %v", err)
		}
		server = grpc.NewServer(
			grpc.Creds(creds),
		)
	} else {
		server = grpc.NewServer()
	}
	reflection.Register(server)

	return &Server{config: cfg, agent: NewDataProviderServiceAgent(service), server: server}, nil
}

func (s *Server) Serve() error {
	pb.RegisterDataProviderServiceServer(s.server, s.agent)

	lis, err := net.Listen(s.config.Network, s.config.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	log.Printf("serving grpc at %v %v", s.config.Network, s.config.Addr)
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
