package gapi

import (
	"context"
	"log"
	"net"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/config"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/generated"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/observability"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/token_maker"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
	// initialize tracer provider
	ctx1 := context.Background()
	trace_cleanup, err := observability.InitTracer(ctx1, s.config)
	if err != nil {
		log.Fatal("---> Could not initialize tracer provider: ", err)
	}
	defer trace_cleanup()

	// initialize meter provider
	ctx2 := context.Background()
	metric_cleanup, err := observability.InitMeterProvider(ctx2, s.config)
	if err != nil {
		log.Fatal("---> Could not initilize meter provider: ", err)
	}
	defer metric_cleanup(context.Background())

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	generated.RegisterAuthGrpcServiceServer(grpcServer, s)

	listener, err := net.Listen("tcp", s.config.GRPCAddress)
	if err != nil {
		log.Fatal("---> Could create listener: ", err)
	}

	log.Println("---> Server starting listening at: ", s.config.GRPCAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Could not listen at address: %s due to Error: %s", s.config.GRPCAddress, err)
	}
}
