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
		log.Printf("[OK] Created: ID=%d, Name=%s, Email=%s",
			createResp.GetId(),
			createResp.GetName(),
			createResp.GetEmail())
	}

	// Test 2: Update user
	log.Println("\n=== Test 2: Update User ===")
	var updatedEmail string
	if createResp != nil {
		updatedEmail = fmt.Sprintf("updated%d@example.com", timestamp)
		updateResp, err := client.UpdateUser(ctx, &userpb.UpdateUserRequest{
			Id:    createResp.GetId(),
			Name:  "Updated User",
			Email: updatedEmail,
		})
		if err != nil {
			handleGRPCError("UpdateUser", err)
		} else {
			log.Printf("[OK] Updated: ID=%d, Name=%s, Email=%s",
				updateResp.GetId(),
				updateResp.GetName(),
				updateResp.GetEmail())
		}
	}

	// Test 3: List users dengan pagination
	log.Println("\n=== Test 3: List Users (limit 5) ===")
	listResp, err := client.ListUsers(ctx, &userpb.ListUsersRequest{
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		handleGRPCError("ListUsers", err)
	} else {
		log.Printf("[OK] Total users: %d", listResp.GetTotal())
		log.Println("Users:")
		for i, user := range listResp.GetUsers() {
			log.Printf("  %d. ID=%d, Name=%s, Email=%s", i+1, user.GetId(), user.GetName(), user.GetEmail())
		}
	}

	// Test 4: Get user yang baru dibuat
	if createResp != nil {
		log.Println("\n=== Test 4: Get Existing User ===")
		getResp, err := client.GetUser(ctx, &userpb.GetUserRequest{
			Id: createResp.GetId(),
		})
		if err != nil {
			handleGRPCError("GetUser", err)
		} else {
			log.Printf("[OK] Got: ID=%d, Name=%s, Email=%s",
				getResp.GetId(),
				getResp.GetName(),
				getResp.GetEmail())
		}
	}

	// Test 5: Get user yang tidak ada
	log.Println("\n=== Test 5: Get Non-Existing User ===")
	_, err = client.GetUser(ctx, &userpb.GetUserRequest{
		Id: 999999,
	})
	if err != nil {
		handleGRPCError("GetUser", err)
	}

	// Test 6: Create user dengan email duplikat
	log.Println("\n=== Test 6: Create User with Duplicate Email ===")
	_, err = client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:  "Another User",
		Email: updatedEmail, // Email yang sudah diupdate di test 2
	})
	if err != nil {
		handleGRPCError("CreateUser", err)
	}

	// Test 7: Update user yang tidak ada
	log.Println("\n=== Test 7: Update Non-Existing User ===")
	_, err = client.UpdateUser(ctx, &userpb.UpdateUserRequest{
		Id:    999999,
		Name:  "Ghost User",
		Email: "ghost@example.com",
	})
	if err != nil {
		handleGRPCError("UpdateUser", err)
	}

	// Test 8: Update user dengan input invalid
	log.Println("\n=== Test 8: Update User with Invalid Input ===")
	if createResp != nil {
		_, err = client.UpdateUser(ctx, &userpb.UpdateUserRequest{
			Id:    createResp.GetId(),
			Name:  "",
			Email: "invalid-email",
		})
		if err != nil {
			handleGRPCError("UpdateUser", err)
		}
	}

	log.Println("\n[DONE] All tests completed!")
}

func handleGRPCError(method string, err error) {
	st, ok := status.FromError(err)
	if !ok {
		log.Printf("[ERROR] %s failed: %v", method, err)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		log.Printf("[WARN] %s: %s (code: NotFound)", method, st.Message())
	case codes.AlreadyExists:
		log.Printf("[WARN] %s: %s (code: AlreadyExists)", method, st.Message())
	case codes.InvalidArgument:
		log.Printf("[WARN] %s: %s (code: InvalidArgument)", method, st.Message())
	case codes.Internal:
		log.Printf("[ERROR] %s: %s (code: Internal)", method, st.Message())
	default:
		log.Printf("[ERROR] %s: %s (code: %s)", method, st.Message(), st.Code())
	}
}
