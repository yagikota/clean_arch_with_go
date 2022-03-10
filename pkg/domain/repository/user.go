package repository

import "22dojo-online/pkg/domain/model"

// 依存性の逆転を適用, 実装はinfra層に任せる
// interfaceにすることでアクセスレベルを制限できる, https://selfnote.work/20201123/programming/how-to-use-interface-in-golang/
// TODO: when transaction
type UserRepository interface {
	CreateUser(record *model.User) error 
	SelectUserByAuthToken(string) (*model.User, error)
	SelectUserByPrimaryKey(string) (*model.User, error)
	UpdateUserByPrimaryKey(user *model.User) error
}
