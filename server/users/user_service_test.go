package users_test

import (
	"context"
	"testing"

	protos "github.com/cndrsdrmn/go-grpc/protos/gen"
	"github.com/cndrsdrmn/go-grpc/server/users"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func setupUserServiceTest(t *testing.T) users.UserServiceInterface {
	defer teardownTest(t)
	repo := setupUserRepositoryTest(t)
	return users.NewUserService(repo)
}

func TestSrvsCreateUser(t *testing.T) {
	srvs := setupUserServiceTest(t)
	ctx := context.Background()
	req := &protos.CreateUserRequest{Name: "John Doe", Email: "john@example.com", Password: "secret"}

	t.Run("success create a new user", func(t *testing.T) {
		res, err := srvs.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), res.User.Id)
		assert.Equal(t, "John Doe", res.User.Name)
		assert.Equal(t, "john@example.com", res.User.Email)
	})

	t.Run("failed create an existing user", func(t *testing.T) {
		_, err := srvs.CreateUser(ctx, req)

		assert.Error(t, err, gorm.ErrDuplicatedKey)
	})
}

func TestSrvsGetUser(t *testing.T) {
	srvs := setupUserServiceTest(t)
	ctx := context.Background()

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	t.Run("find an existing user", func(t *testing.T) {
		res, err := srvs.GetUser(ctx, &protos.GetUserRequest{Id: uint64(user.ID)})

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), res.User.Id)
		assert.Equal(t, "John Doe", res.User.Name)
		assert.Equal(t, "john@example.com", res.User.Email)
	})

	t.Run("find a non-existing user", func(t *testing.T) {
		_, err := srvs.GetUser(ctx, &protos.GetUserRequest{Id: 999})

		assert.Error(t, err, gorm.ErrRecordNotFound)
	})
}

func TestSrvsAllUsers(t *testing.T) {
	srvs := setupUserServiceTest(t)
	ctx := context.Background()

	factoryUserCreate(&users.User{Name: "Charlie", Email: "charlie@example.com", Password: "secret"})
	factoryUserCreate(&users.User{Name: "David", Email: "david@example.com", Password: "password"})

	res, err := srvs.AllUsers(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.Len(t, res.Users, 2)
}

func TestSrvsUpdateUser(t *testing.T) {
	srvs := setupUserServiceTest(t)
	ctx := context.Background()

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	res, err := srvs.UpdateUser(ctx, &protos.UpdateUserRequest{Id: uint64(user.ID), Name: "Charlie"})

	assert.NoError(t, err)
	assert.Equal(t, "Charlie", res.User.Name)
	assert.Equal(t, "john@example.com", res.User.Email)
}

func TestSrvsDeleteUser(t *testing.T) {
	srvs := setupUserServiceTest(t)
	ctx := context.Background()

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	t.Run("can delete an existing user", func(t *testing.T) {
		res, err := srvs.DeleteUser(ctx, &protos.DeleteUserRequest{Id: uint64(user.ID)})

		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("cannot delete a non-existing user", func(t *testing.T) {
		res, err := srvs.DeleteUser(ctx, &protos.DeleteUserRequest{Id: uint64(user.ID)})

		assert.Error(t, err, gorm.ErrRecordNotFound)
		assert.False(t, res.Success)
	})
}
