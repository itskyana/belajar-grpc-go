package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/itskyana/belajar-grpc-go/gen/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Creating user...")

	timestamp := time.Now().Unix()
	email := fmt.Sprintf("user%d@example.com", timestamp)
	name := fmt.Sprintf("User %d", timestamp)

	log.Printf("Generated email: %s", email)
	log.Printf("Generated name: %s", name)
	createResp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("User created: ID=%d, Name=%s, Email=%s", createResp.GetId(), createResp.GetName(), createResp.GetEmail())

	log.Println("Retrieving user...")
	getResp, err := client.GetUser(ctx, &userpb.GetUserRequest{
		Id: createResp.GetId(),
	})
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	log.Printf("User retrieved: ID=%d, Name=%s, Email=%s", getResp.GetId(), getResp.GetName(), getResp.GetEmail())

	log.Println("All operations completed successfully.")
}
