package service

import (
	"22dojo-online/pkg/domain/model"
	"22dojo-online/pkg/domain/repository"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository: userRepository}
}

func (us *userService) CreateUser(record *model.User) error {
	err := us.userRepository.CreateUser(record)
	return err
}

func (us *userService) SelectUserByAuthToken(authToken string) (*model.User, error) {
	user, err := us.userRepository.SelectUserByAuthToken(authToken)
	return user, err
}

func (us *userService) SelectUserByPrimaryKey(userID string) (*model.User, error) {
	user, err := us.userRepository.SelectUserByPrimaryKey(userID)
	return user, err
}

func (us *userService) UpdateUserByPrimaryKey(record *model.User) error {
	err := us.userRepository.UpdateUserByPrimaryKey(record)
	return err
}
