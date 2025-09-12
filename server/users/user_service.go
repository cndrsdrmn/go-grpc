package users

import (
	"context"

	pb "github.com/cndrsdrmn/go-grpc/protos/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceInterface interface {
	pb.UserServiceServer
}

type userService struct {
	pb.UnimplementedUserServiceServer
	repo UserRepositoryInterface
}

func (srvs *userService) AllUsers(context.Context, *emptypb.Empty) (*pb.AllUsersResponse, error) {
	users, err := srvs.repo.AllUser()
	if err != nil {
		return nil, err
	}

	var res []*pb.User
	for _, u := range users {
		res = append(res, u.ToProtoUserResponse().User)
	}

	return &pb.AllUsersResponse{Users: res}, nil
}

func (srvs *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := srvs.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user.ToProtoUserResponse(), nil
}

func (srvs *userService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := srvs.repo.DeleteUser(uint(req.Id))
	if err != nil {
		return &pb.DeleteUserResponse{Success: false}, err
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

func (srvs *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := srvs.repo.FindUser(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return user.ToProtoUserResponse(), nil
}

func (srvs *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := srvs.repo.FindUser(uint(req.Id))
	if err != nil {
		return nil, err
	}

	user.Name = req.Name

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Password != nil {
		user.Password = *req.Password
	}

	if err := srvs.repo.UpdateUser(user.ID, user); err != nil {
		return nil, err
	}

	return user.ToProtoUserResponse(), nil
}

func NewUserService(repo UserRepositoryInterface) UserServiceInterface {
	return &userService{repo: repo}
}
