package db

// pkg/infrastructure/sql/conn.goを直接参照させないDIP（依存関係逆転の原則）

type SQLHandler interface {
	Exec(string, ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(string, ...interface{}) Row
}

// https://pkg.go.dev/database/sql#Result
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Row interface {
	Err() error
	Scan(...interface{}) error
}

type Rows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}
