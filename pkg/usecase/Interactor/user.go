package interactor

import (
	"errors"

	"github.com/google/uuid"

	"22dojo-online/pkg/domain/model"
	"22dojo-online/pkg/domain/service"
	inputdata "22dojo-online/pkg/usecase/input_data"
)

type UserInteractor interface {
	CreateUser(inputdata.UserCreateRequest) (string, error)
	SelectUserByAuthToken(string) (*model.User, error)
	SelectUserByPrimaryKey(string) (*model.User, error)
	UpdateUserByPrimaryKey(inputdata.UserUpdateRequest, string) error
}

type userInteractor struct {
	userService service.UserService
}

func NewUserInteractor(userService service.UserService) UserInteractor {
	return &userInteractor{
		userService: userService,
	}
}

func (ui *userInteractor) CreateUser(requestBody inputdata.UserCreateRequest) (string, error) {
	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// UUIDで認証トークンを生成する
	authToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// データベースにユーザデータを登録する
	user := &model.User{
		ID:        userID.String(),
		AuthToken: authToken.String(),
		Name:      requestBody.Name,
		HighScore: 0,
		Coin:      0,
	}
	if err := ui.userService.CreateUser(user); err != nil {
		return "", err
	}

	// 生成した認証トークンを返却
	return user.AuthToken, nil
}

func (ui *userInteractor) SelectUserByPrimaryKey(userID string) (*model.User, error) {
	return ui.userService.SelectUserByPrimaryKey(userID)
}

func (ui *userInteractor) SelectUserByAuthToken(authToken string) (*model.User, error) {
	return ui.userService.SelectUserByAuthToken(authToken)
}

func (ui *userInteractor) UpdateUserByPrimaryKey(requestBody inputdata.UserUpdateRequest, userID string) error {
	user, err := ui.userService.SelectUserByPrimaryKey(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Name = requestBody.Name
	return ui.userService.UpdateUserByPrimaryKey(user)
}
