package db

// pkg/infrastructure/sql/conn.goを直接参照させないDIP（依存関係逆転の原則）

type SqlHandler interface {
	Exec(string, ...interface{}) (Result, error)
	QueryRow(string, ...interface{}) Row
}

// https://pkg.go.dev/database/sql#Result
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Row interface {
	Scan(...interface{}) error
	Err() error
}
