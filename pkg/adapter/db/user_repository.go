package db

import (
	"database/sql"
	"log"

	"22dojo-online/pkg/domain/model"

	"github.com/pkg/errors"
)

type UserRepository struct {
	SQLHandler
}

func NewUserRepository(sqlHandler SQLHandler) *UserRepository {
	return &UserRepository{
		SQLHandler: sqlHandler,
	}
}

// CreateUser データベースをレコードを登録する
func (repo *UserRepository) CreateUser(record *model.User) error {
	// TODO: usersテーブルへのレコードの登録を行うSQLを入力する
	_, err := repo.Exec(
		"INSERT INTO user VALUES (?, ?, ?, ?, ?)",
		record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
	return errors.Wrap(err, "[Error] failed to execute query, CreateUser")
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func (repo *UserRepository) SelectUserByAuthToken(authToken string) (*model.User, error) {
	// TODO: auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
	row := repo.QueryRow("SELECT id, auth_token, name, high_score, coin FROM user WHERE auth_token=?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func (repo *UserRepository) SelectUserByPrimaryKey(userID string) (*model.User, error) {
	// TODO: idを条件にSELECTを行うSQLを第1引数に入力する
	row := repo.QueryRow("SELECT id, auth_token, name, high_score, coin FROM user WHERE id=?", userID)
	return convertToUser(row)
}

// 排他制御
func (repo *UserRepository) SelectUserByPrimaryKeyWithLock(tx *sql.Tx, userID string) (*model.User, error) {
	row := tx.QueryRow("SELECT id, auth_token, name, high_score, coin FROM user WHERE id=? FOR UPDATE", userID)
	return convertToUser(row)
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func (repo *UserRepository) UpdateUserByPrimaryKey(record *model.User) error {
	// TODO: idを条件に指定した値でnameカラムの値を更新するSQLを入力する
	_, err := repo.Exec(
		"UPDATE user SET name=? WHERE id=?",
		record.Name, record.ID)
	return errors.Wrap(err, "[Error] failed to execute query, UpdateUserByPrimaryKey")
}

func (repo *UserRepository) UpdateUserByPrimaryKeyWithTx(tx *sql.Tx, record *model.User) error {
	_, err := tx.Exec(
		"UPDATE user SET id=?, auth_token=?, name=?, high_score=?, coin=? WHERE id=?",
		record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin, record.ID)
	return errors.Wrap(err, "[Error] failed to execute query, UpdateUserByPrimaryKeyWithTx")
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row Row) (*model.User, error) {
	user := model.User{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, errors.Wrap(err, "[Error] failed to scan row, convertToUser")
	}
	return &user, nil
}
