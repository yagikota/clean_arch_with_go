package inputboundary

import (
	"22dojo-online/pkg/domain/model"
	inputdata "22dojo-online/pkg/usecase/input_data"
)

type UserCreateRequestBoundary interface {
	CreateUser(inputdata.UserCreateRequest) (string, error)
}

type UserGetRequestBoundary interface {
	GetUser(string) (*model.User, error)
	SelectUserByAuthToken(string) (*model.User, error)
}

type UserUpdataRequestBoundary interface {
	GetUser(string) (*model.User, error)
	UpdateUser(inputdata.UserUpdateRequest, string) error
	SelectUserByAuthToken(string) (*model.User, error)
}
