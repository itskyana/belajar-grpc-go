package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/itskyana/belajar-grpc-go/gen/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test 1: Create user valid
	log.Println("=== Test 1: Create Valid User ===")
	timestamp := time.Now().Unix()
	email := fmt.Sprintf("user%d@example.com", timestamp)
	name := fmt.Sprintf("User %d", timestamp)

	createResp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		handleGRPCError("CreateUser", err)
	} else {
		log.Printf("Created: ID=%d, Name=%s, Email=%s",
			createResp.GetId(),
			createResp.GetName(),
			createResp.GetEmail())
	}

	// Test 2: Get user yang baru dibuat
	if createResp != nil {
		log.Println("\n=== Test 2: Get Existing User ===")
		getResp, err := client.GetUser(ctx, &userpb.GetUserRequest{
			Id: createResp.GetId(),
		})
		if err != nil {
			handleGRPCError("GetUser", err)
		} else {
			log.Printf("Got: ID=%d, Name=%s, Email=%s",
				getResp.GetId(),
				getResp.GetName(),
				getResp.GetEmail())
		}
	}

	// Test 3: Get user yang tidak ada
	log.Println("\n=== Test 3: Get Non-Existing User ===")
	_, err = client.GetUser(ctx, &userpb.GetUserRequest{
		Id: 999999,
	})
	if err != nil {
		handleGRPCError("GetUser", err)
	}

	// Test 4: Create user dengan email duplikat
	log.Println("\n=== Test 4: Create User with Duplicate Email ===")
	_, err = client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  "Another User",
		Email: email, // Email yang sama
	})
	if err != nil {
		handleGRPCError("CreateUser", err)
	}

	// Test 5: Create user dengan input invalid
	log.Println("\n=== Test 5: Create User with Invalid Input ===")
	_, err = client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  "",
		Email: "invalid-email",
	})
	if err != nil {
		handleGRPCError("CreateUser", err)
	}

	log.Println("\nAll tests completed!")
}

func handleGRPCError(method string, err error) {
	st, ok := status.FromError(err)
	if !ok {
		log.Printf("%s failed: %v", method, err)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		log.Printf("%s: %s (code: NotFound)", method, st.Message())
	case codes.AlreadyExists:
		log.Printf("%s: %s (code: AlreadyExists)", method, st.Message())
	case codes.InvalidArgument:
		log.Printf("%s: %s (code: InvalidArgument)", method, st.Message())
	case codes.Internal:
		log.Printf("%s: %s (code: Internal)", method, st.Message())
	default:
		log.Printf("%s: %s (code: %s)", method, st.Message(), st.Code())
	}
}
