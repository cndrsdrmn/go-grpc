package main

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	protos "github.com/cndrsdrmn/go-grpc/protos/gen"
	"github.com/cndrsdrmn/go-grpc/server/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var bufSize = 1024 * 1024
var lis *bufconn.Listener
var testDB *gorm.DB

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	db.AutoMigrate(&users.User{})

	lis = bufconn.Listen(bufSize)

	server := NewGRPCServer(db)
	testDB = db

	go func() {
		if err := server.Serve(lis); err != nil {
			assert.Fail(nil, "Server failed to serve: %v", err)
		}
	}()

	code := m.Run()

	server.Stop()

	os.Exit(code)
}

func dialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func factoryUserCreate(user *users.User) error {
	return testDB.Create(user).Error
}

func setupTestServer(t *testing.T) (*grpc.ClientConn, func()) {
	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	assert.NoError(t, err)

	cleanup := func() {
		conn.Close()
		teardownTest(t)
	}

	return conn, cleanup
}

func teardownTest(t *testing.T) {
	if err := testDB.Exec("DELETE FROM users").Error; err != nil {
		t.Fatalf("failed to clear users table: %v", err)
	}

	if err := testDB.Exec("DELETE FROM sqlite_sequence WHERE name='users'").Error; err != nil {
		t.Fatalf("failed to reset autoincrement sequence: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	conn, cleanup := setupTestServer(t)
	defer cleanup()

	client := protos.NewUserServiceClient(conn)

	ctx := context.Background()
	req := &protos.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret",
	}

	t.Run("success create a new user", func(t *testing.T) {
		res, err := client.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), res.User.Id)
		assert.Equal(t, "John Doe", res.User.Name)
		assert.Equal(t, "john@example.com", res.User.Email)
	})

	t.Run("failed create an existing user", func(t *testing.T) {
		_, err := client.CreateUser(ctx, req)

		assert.Error(t, err, gorm.ErrDuplicatedKey)
	})
}

func TestGetUser(t *testing.T) {
	conn, cleanup := setupTestServer(t)
	defer cleanup()

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	client := protos.NewUserServiceClient(conn)

	ctx := context.Background()

	t.Run("find an existing user", func(t *testing.T) {
		res, err := client.GetUser(ctx, &protos.GetUserRequest{Id: uint64(user.ID)})

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), res.User.Id)
		assert.Equal(t, "John Doe", res.User.Name)
		assert.Equal(t, "john@example.com", res.User.Email)
	})

	t.Run("find a non-existing user", func(t *testing.T) {
		_, err := client.GetUser(ctx, &protos.GetUserRequest{Id: 999})

		assert.Error(t, err, gorm.ErrRecordNotFound)
	})
}

func TestAllUsers(t *testing.T) {
	conn, cleanup := setupTestServer(t)
	defer cleanup()

	client := protos.NewUserServiceClient(conn)

	ctx := context.Background()

	factoryUserCreate(&users.User{Name: "Charlie", Email: "charlie@example.com", Password: "secret"})
	factoryUserCreate(&users.User{Name: "David", Email: "david@example.com", Password: "password"})

	res, err := client.AllUsers(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.Len(t, res.Users, 2)
}

func TestUpdateUser(t *testing.T) {
	conn, cleanup := setupTestServer(t)
	defer cleanup()

	client := protos.NewUserServiceClient(conn)

	ctx := context.Background()

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	res, err := client.UpdateUser(ctx, &protos.UpdateUserRequest{Id: uint64(user.ID), Name: "Charlie"})

	assert.NoError(t, err)
	assert.Equal(t, "Charlie", res.User.Name)
	assert.Equal(t, "john@example.com", res.User.Email)
}

func TestDeleteUser(t *testing.T) {
	conn, cleanup := setupTestServer(t)
	defer cleanup()

	client := protos.NewUserServiceClient(conn)

	ctx := context.Background()

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	res, err := client.DeleteUser(ctx, &protos.DeleteUserRequest{Id: uint64(user.ID)})

	assert.NoError(t, err)
	assert.True(t, res.Success)
}
