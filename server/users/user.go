package users

import (
	pb "github.com/cndrsdrmn/go-grpc/protos/gen"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	return user.hashPassword(tx)
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	if !tx.Statement.Changed("Password") {
		return nil
	}

	return user.hashPassword(tx)
}

func (user *User) hashPassword(tx *gorm.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Password", string(hashed))
	return nil
}

func (user User) ToProtoUserResponse() *pb.UserResponse {
	return &pb.UserResponse{
		User: &pb.User{
			Id:    uint64(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}
}
