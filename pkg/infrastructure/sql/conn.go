package db

// DB接続には、外部パッケージを使用しているので、infrastructure層に定義し外側のルールを内側に持ち込まないようにします。

import (
	"22dojo-online/pkg/adapter/db"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	// blank import for MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Driver名
const driverName = "mysql"

// Conn 各repositoryで利用するDB接続(Connection)情報
// var Conn *sql.DB
type SqlHandler struct {
	Conn *sql.DB
}

func NewSqlHandler() *SqlHandler {
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

	return &SqlHandler{
		Conn: conn,
	}
}

	func (handler *SqlHandler) Exec(query string, args ...interface{}) (db.Result, error) {
		res := SqlResult{}
		result, err := handler.Conn.Exec(query, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

func (handler *SqlHandler) QueryRow(query string, args ...interface{}) db.Row {
    return handler.Conn.QueryRow(query, args...)
}

type SqlResult struct {
	Result sql.Result
}

func (r SqlResult) LastInsertId() (int64, error) {
    return r.Result.LastInsertId()
}

func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type SqlRow struct {
	Row *sql.Row
}

func (r SqlRow) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

func (r SqlRow) Err() error {
	return r.Row.Err()
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
