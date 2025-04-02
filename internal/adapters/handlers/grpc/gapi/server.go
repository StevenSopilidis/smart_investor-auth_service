package gapi

import (
	"log"
	"net"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/config"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/generated"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/token_maker"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services"
	"google.golang.org/grpc"
)

type Server struct {
	generated.UnimplementedAuthGrpcServiceServer
	auth_service *services.AuthService[token_maker.Claims]
	config       config.Config
}

func NewServer(
	auth_service *services.AuthService[token_maker.Claims],
	config config.Config,
) *Server {
	return &Server{
		auth_service: auth_service,
		config:       config,
	}
}

func (s *Server) Run() {
	grpc := grpc.NewServer()
	generated.RegisterAuthGrpcServiceServer(grpc, s)

	listener, err := net.Listen("tcp", s.config.GRPCAddress)
	if err != nil {
		log.Fatal("---> Could create listener: ", err)
	}

	log.Println("---> Server starting listening at: ", s.config.GRPCAddress)
	err = grpc.Serve(listener)
	if err != nil {
		log.Fatalf("Could not listen at address: %s due to Error: %s", s.config.GRPCAddress, err)
	}
}
