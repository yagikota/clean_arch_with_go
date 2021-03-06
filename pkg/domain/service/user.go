package service

// サービスは以下の用途で利用します。
// 	ドメインオブジェクトに責務を持たせるものではないケース
// 	データ整合性を保つために複数のドメインモデルを操作するケース
// 前提としてサービスはステートレス(入力の内容によってのみ出力が決定される)である必要があります。
// serviceが依存するのはrepository(IF)とmodelだけ
// 今回の場合は、わざわざserviceを挟む必要はない。

import (
	"22dojo-online/pkg/domain/model"
	"22dojo-online/pkg/domain/repository"
)

type UserService interface {
	CreateUser(record *model.User) error
	SelectUserByAuthToken(userID string) (*model.User, error)
	SelectUserByPrimaryKey(userID string) (*model.User, error)
	UpdateUserByPrimaryKey(user *model.User) error
}

type userService struct {
	Repository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		Repository: userRepository,
	}
}

func (us *userService) CreateUser(record *model.User) error {
	return us.Repository.CreateUser(record)
}

func (us *userService) SelectUserByAuthToken(authToken string) (*model.User, error) {
	user, err := us.Repository.SelectUserByAuthToken(authToken)
	return user, err
}

func (us *userService) SelectUserByPrimaryKey(userID string) (*model.User, error) {
	user, err := us.Repository.SelectUserByPrimaryKey(userID)
	return user, err
}

func (us *userService) UpdateUserByPrimaryKey(record *model.User) error {
	return us.Repository.UpdateUserByPrimaryKey(record)
}
