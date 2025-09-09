package main

import (
	"context"
	"fmt"
	"log"
	"time"

	protos "github.com/cndrsdrmn/go-grpc/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Cannot connect to: %v", err)
	}
	defer conn.Close()

	client := protos.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Create User
	createResp, err := client.CreateUser(ctx, &protos.CreateUserRequest{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Println("Created:", createResp.User)

	// 2. Get User
	getResp, err := client.GetUser(ctx, &protos.GetUserRequest{Id: createResp.User.Id})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	fmt.Println("Fetched:", getResp.User)

	// 3. Update User
	email := "alice.new@example.com"
	updateResp, err := client.UpdateUser(ctx, &protos.UpdateUserRequest{
		Id:    createResp.User.Id,
		Name:  "Alice Updated",
		Email: &email,
	})
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	fmt.Println("Updated:", updateResp.User)

	// 4. List Users
	listResp, err := client.AllUsers(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	fmt.Println("List Users:", listResp.Users)

	// 5. Delete User
	delResp, err := client.DeleteUser(ctx, &protos.DeleteUserRequest{Id: createResp.User.Id})
	if err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	fmt.Println("Deleted:", delResp.Success)
}
