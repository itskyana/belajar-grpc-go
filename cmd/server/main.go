package main

import (
	"log"
	"net"

	"github.com/itskyana/belajar-grpc-go/gen/userpb"
	"github.com/itskyana/belajar-grpc-go/internal/config"
	"github.com/itskyana/belajar-grpc-go/internal/database"
	grpcdelivery "github.com/itskyana/belajar-grpc-go/internal/delivery/grpc"
	"github.com/itskyana/belajar-grpc-go/internal/repository/postgres"
	"github.com/itskyana/belajar-grpc-go/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("failed to create postgres pool: %v", err)
	}

	defer pool.Close()
	log.Println("Connected to PostgreSQL database")

	userRepo := postgres.NewUserRepository(pool)

	userUseCase := usecase.NewUserUseCase(userRepo)

	userHandler := grpcdelivery.NewUserHandler(*userUseCase)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server is running on port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
