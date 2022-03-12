package model

import (
	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/db"

	"github.com/pkg/errors"
)

type Rank struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

// TODO: テーブルにインデックス貼る データ数増えるとかなり重い検索になってしまうので、対応した方が良さそう。（発展課題以降）
// TODO: 共通化やテスト、コードの自動生成等を考慮するとDB検索処理と検索結果を加工するロジックは分離しておいた方が扱いやすいかも
func GetRanking(start int) ([]*Rank, error) {
	limit := constant.RankingGetLimit
	rows, err := db.Conn.Query("SELECT id, name, high_score FROM user ORDER BY high_score DESC, id ASC LIMIT ? OFFSET ?", limit, start-1)
	if err != nil {
		return nil, errors.Wrap(err, "[Error] failed to execute query, GetRanking")
	}
	defer rows.Close()

	var rankingList []*Rank
	currRank := start
	for rows.Next() {
		rank := Rank{Rank: currRank}
		err := rows.Scan(&rank.UserID, &rank.UserName, &rank.Score)
		if err != nil {
			return nil, errors.Wrap(err, "[Error] failed to scan rows, GetRanking")
		}
		rankingList = append(rankingList, &rank)
		currRank++
	}

	return rankingList, nil
}

func CountRecord() (int, error) {
	row := db.Conn.QueryRow("SELECT COUNT(id) FROM user")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, errors.Wrap(err, "[Error] failed tofi scan row, CountRecord")
	}
	return count, nil
}
