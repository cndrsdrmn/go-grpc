package users_test

import (
	"testing"

	"github.com/cndrsdrmn/go-grpc/server/users"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupUserRepositoryTest(t *testing.T) users.UserRepositoryInterface {
	defer teardownTest(t)
	return users.NewUserRepository(testDB)
}

func TestRepoCreateUser(t *testing.T) {
	repo := setupUserRepositoryTest(t)
	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}

	t.Run("success create a new user", func(t *testing.T) {
		err := repo.CreateUser(user)

		assert.NoError(t, err)
		assert.NotEqual(t, uint(0), user.ID)
	})

	t.Run("failed create an existing user", func(t *testing.T) {
		err := repo.CreateUser(user)

		assert.Error(t, err, gorm.ErrDuplicatedKey)
	})
}

func TestRepoFindUser(t *testing.T) {
	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	repo := setupUserRepositoryTest(t)
	factoryUserCreate(user)

	t.Run("find an existing user", func(t *testing.T) {
		founded, err := repo.FindUser(user.ID)

		assert.NoError(t, err)
		assert.Equal(t, user.Name, founded.Name)
		assert.Equal(t, user.Email, founded.Email)

		err = bcrypt.CompareHashAndPassword([]byte(founded.Password), []byte("secret"))
		assert.NoError(t, err)
	})

	t.Run("find a non-existing user", func(t *testing.T) {
		_, err := repo.FindUser(999)

		assert.Error(t, err, gorm.ErrRecordNotFound)
	})
}

func TestRepoAllUsers(t *testing.T) {
	repo := setupUserRepositoryTest(t)
	factoryUserCreate(&users.User{Name: "Charlie", Email: "charlie@example.com", Password: "secret"})
	factoryUserCreate(&users.User{Name: "David", Email: "david@example.com", Password: "password"})

	users, err := repo.AllUser()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestRepoUpdateUser(t *testing.T) {
	repo := setupUserRepositoryTest(t)
	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	updated := &users.User{Name: "Lorem Ipsum", Email: "lorem@example.com", Password: "supersecret"}

	err := repo.UpdateUser(user.ID, updated)

	assert.NoError(t, err)

	assert.NotEqual(t, user.Email, updated.Email)
	assert.NotEqual(t, user.Name, updated.Name)
	assert.NotEqual(t, user.Password, updated.Password)
}

func TestRepoDeleteUser(t *testing.T) {
	repo := setupUserRepositoryTest(t)
	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	t.Run("can delete an existing user", func(t *testing.T) {
		err := repo.DeleteUser(user.ID)

		assert.NoError(t, err)
	})

	t.Run("cannot delete a non-existing user", func(t *testing.T) {
		err := repo.DeleteUser(999)

		assert.Error(t, err, gorm.ErrRecordNotFound)
	})
}
