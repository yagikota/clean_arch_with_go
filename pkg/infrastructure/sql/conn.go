package infrasql

// DB接続には、外部パッケージを使用しているので、infrastructure層に定義し外側のルールを内側に持ち込まないようにします。

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	// blank import for MySQL driver
	_ "github.com/go-sql-driver/mysql"

	"22dojo-online/pkg/adapter/db"
)

// Driver名
const driverName = "mysql"

// Conn 各repositoryで利用するDB接続(Connection)情報
// var Conn *sql.DB
type SQLHandler struct {
	Conn *sql.DB
}

func NewSQLHandler() *SQLHandler {
	/* ===== データベースへ接続する. ===== */
	// ユーザ
	user := os.Getenv("MYSQL_USER")
	// パスワード
	password := os.Getenv("MYSQL_PASSWORD")
	// 接続先ホスト
	host := os.Getenv("MYSQL_HOST")
	// 接続先ポート
	port := os.Getenv("MYSQL_PORT")
	// 接続先データベース
	database := os.Getenv("MYSQL_DATABASE")

	// 接続情報は以下のように指定する.
	// user:password@tcp(host:port)/database
	conn, err := sql.Open(driverName,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatalf("can't connect to mysql server. "+
			"MYSQL_USER=%s, "+
			"MYSQL_PASSWORD=%s, "+
			"MYSQL_HOST=%s, "+
			"MYSQL_PORT=%s, "+
			"MYSQL_DATABASE=%s, "+
			"error=%+v",
			user, password, host, port, database, err)
	}

	return &SQLHandler{
		Conn: conn,
	}
}

func (handler *SQLHandler) Exec(query string, args ...interface{}) (db.Result, error) {
	res := SQLResult{}
	result, err := handler.Conn.Exec(query, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

//nolint: rowserrcheck // this is why
func (handler *SQLHandler) Query(query string, args ...interface{}) (db.Rows, error) {
	return handler.Conn.Query(query, args...)
}

func (handler *SQLHandler) QueryRow(query string, args ...interface{}) db.Row {
	return handler.Conn.QueryRow(query, args...)
}

type SQLResult struct {
	Result sql.Result
}

func (r SQLResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SQLResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type SQLRow struct {
	Row *sql.Row
}

func (r SQLRow) Err() error {
	return r.Row.Err()
}

func (r SQLRow) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

type SQLRows struct {
	Rows *sql.Rows
}

func (r SQLRows) Close() error {
	return r.Rows.Close()
}

func (r SQLRows) Err() error {
	return r.Rows.Err()
}

func (r SQLRows) Next() bool {
	return r.Rows.Next()
}

func (r SQLRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

// トランザクションのWrapper
// TODO: tx.Rollback()のerror handling
// TODO: logに出すerror内容もっと詳しい方がいいか？
func ExecWithTx(ctx context.Context, conn *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println(rollbackErr)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println(rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Println(commitErr)
			}
		}
	}()

	return txFunc(tx)
}
