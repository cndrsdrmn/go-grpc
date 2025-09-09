package users

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	AllUser() ([]User, error)
	CreateUser(user *User) error
	FindUser(id uint) (*User, error)
	UpdateUser(id uint, user *User) error
	DeleteUser(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) AllUser() ([]User, error) {
	var users []User
	err := repo.db.Find(&users).Error
	return users, err
}

func (repo *userRepository) CreateUser(user *User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) FindUser(id uint) (*User, error) {
	var user User
	err := repo.db.First(&user, id).Error
	return &user, err
}

func (repo *userRepository) UpdateUser(id uint, user *User) error {
	updates := make(map[string]interface{})
	if user.Name != "" {
		updates["Name"] = user.Name
	}
	if user.Email != "" {
		updates["Email"] = user.Email
	}
	if user.Password != "" {
		updates["Password"] = user.Password
	}
	if len(updates) == 0 {
		return errors.New("no fields provided to update")
	}

	res := repo.db.Model(&User{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	updated, err := repo.FindUser(id)
	if err != nil {
		return err
	}

	*user = *updated

	return nil
}

func (repo *userRepository) DeleteUser(id uint) error {
	res := repo.db.Unscoped().Delete(&User{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db}
}
