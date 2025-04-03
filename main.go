package main

import (
	"log"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/config"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/gapi"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/password_verification"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/session_manager"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/token_maker"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/user_service_client/grpc_client"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("---> Error could not load config: ", err)
	}

	token_maker, err := token_maker.NewPasetoTokenMaker(config.SymmetricKey)
	if err != nil {
		log.Fatal("---> Error could not create token_maker: ", err)
	}

	session_manager, err := session_manager.NewRedisSessionManager(config, token_maker)
	if err != nil {
		log.Fatal("---> Error could not create session manager: ", err)
	}

	password_verification := password_verification.NewBcryptPasswordHashService()

	user_client, err := grpc_client.NewGRPCUserServiceClient(config.UserServiceAddr)
	if err != nil {
		log.Fatal("---> Could not create grpc_user_service_client: ", err)
	}

	auth_service := services.NewAuthService(session_manager, password_verification, user_client)
	server := gapi.NewServer(auth_service, config)
	server.Run()
}
