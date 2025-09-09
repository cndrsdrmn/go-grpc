package users_test

import (
	"log"
	"os"
	"testing"

	"github.com/cndrsdrmn/go-grpc/server/users"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func factoryUserCreate(user *users.User) error {
	return testDB.Create(user).Error
}

func teardownTest(t *testing.T) {
	if err := testDB.Exec("DELETE FROM users").Error; err != nil {
		t.Fatalf("failed to clear users table: %v", err)
	}

	if err := testDB.Exec("DELETE FROM sqlite_sequence WHERE name='users'").Error; err != nil {
		t.Fatalf("failed to reset autoincrement sequence: %v", err)
	}
}

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&users.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	testDB = db

	code := m.Run()

	ins, err := testDB.DB()
	if err != nil {
		log.Fatalf("failed to get DB instance: %v", err)
	}
	defer ins.Close()

	os.Exit(code)
}

func TestStructBeforeCreate(t *testing.T) {
	defer teardownTest(t)

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}

	err := factoryUserCreate(user)
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret"))
	assert.NoError(t, err)
}

func TestStructBeforeUpdate(t *testing.T) {
	defer teardownTest(t)

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	testDB.Model(&users.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"Email": "charlie@example.com",
	})

	var updated users.User
	err := testDB.First(&updated, user.ID).Error
	assert.NoError(t, err)

	assert.Equal(t, "charlie@example.com", updated.Email)
	assert.Equal(t, user.Name, updated.Name)
	assert.Equal(t, user.Password, updated.Password)
	assert.NotEqual(t, user.Email, updated.Email)
}

func TestStructToProtoUserResponse(t *testing.T) {
	defer teardownTest(t)

	user := &users.User{Name: "John Doe", Email: "john@example.com", Password: "secret"}
	factoryUserCreate(user)

	resp := user.ToProtoUserResponse()

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, uint64(1), resp.User.Id)
	assert.Equal(t, "John Doe", resp.User.Name)
	assert.Equal(t, "john@example.com", resp.User.Email)
}
