package model

import (
	"22dojo-online/pkg/db"
)

type GachaProb struct {
	CollectionID string `json:"collectionID"`
	Ratio        int32  `json:"ratio"`
}

// gacha_probabilityテーブルからid, ratioを取得
func GetAllGachaProbs() ([]*GachaProb, error) {
	rows, err := db.Conn.Query("SELECT collection_item_id, ratio FROM gacha_probability")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gachaProbList []*GachaProb
	for rows.Next() {
		var gachaProb GachaProb
		err = rows.Scan(&gachaProb.CollectionID, &gachaProb.Ratio)
		if err != nil {
			return nil, err
		}
		gachaProbList = append(gachaProbList, &gachaProb)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gachaProbList, nil
}
