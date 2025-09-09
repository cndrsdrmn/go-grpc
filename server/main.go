package main

import (
	"log"
	"net"

	protos "github.com/cndrsdrmn/go-grpc/protos/gen"
	"github.com/cndrsdrmn/go-grpc/server/users"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGRPCServer(db *gorm.DB) *grpc.Server {
	repo := users.NewUserRepository(db)
	srvs := users.NewUserService(repo)

	server := grpc.NewServer()
	protos.RegisterUserServiceServer(server, srvs)
	return server
}

func main() {
	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	db.AutoMigrate(&users.User{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := NewGRPCServer(db)

	log.Println("Server running at :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
