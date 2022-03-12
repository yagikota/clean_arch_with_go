package model

import (
	"database/sql"
	"strings"

	"22dojo-online/pkg/db"
	"22dojo-online/pkg/utils"

	"github.com/pkg/errors"
)

type Collection struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

// ユーザーが所持しているコレクションID一覧を取得
func GetUserCollectionsIDList(userID string) ([]string, error) {
	// 所持アイテムのID取得
	rows, err := db.Conn.Query("SELECT collection_item_id FROM user_collection_item WHERE user_id=?", userID)
	if err != nil {
		return nil, errors.Wrap(err, "[Error] failed to execute query, GetUserCollectionsIDList")
	}
	defer rows.Close()

	var UserCollectionIDList []string
	for rows.Next() {
		var collectionID string
		err = rows.Scan(&collectionID)
		if err != nil {
			return nil, errors.Wrap(err, "[Error] failed to scan rows")
		}
		UserCollectionIDList = append(UserCollectionIDList, collectionID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UserCollectionIDList, nil
}

// DBから取ってきたコレクション
func GetCollectionItems() ([]*Collection, error) {
	rows, err := db.Conn.Query("SELECT id, name rarity FROM collection_item")
	if err != nil {
		return nil, errors.Wrap(err, "[Error] failed to execute query")
		// return nil, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		var collection Collection
		err = rows.Scan(&collection.CollectionID, &collection.Name, &collection.Rarity)
		if err != nil {
			return nil, errors.Wrap(err, "[Error] failed to scan rows")
		}
		collections = append(collections, &collection)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}

// ガチャで出現したコレクション
func GetCollectionItemsByID(idList []string) ([]*Collection, error) {
	valueStrings := make([]string, 0, len(idList))
	valueArgs := make([]interface{}, 0, len(idList))
	for _, id := range idList {
		valueStrings = append(valueStrings, "?")
		valueArgs = append(valueArgs, id)
	}

	query := "SELECT id, name, rarity FROM collection_item WHERE id IN (" + strings.Join(valueStrings, ",") + ")"
	rows, err := db.Conn.Query(query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		var collection Collection
		err = rows.Scan(&collection.CollectionID, &collection.Name, &collection.Rarity)
		if err != nil {
			return nil, err
		}
		collections = append(collections, &collection)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}

func BulkCreateUserCollectionItemWithTx(tx *sql.Tx, userID string, idList []string) error {
	valueStrings := make([]string, 0, len(idList))
	valueArgs := make([]interface{}, 0, len(idList)*2)
	for _, ID := range idList {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, userID, ID)
	}

	query := "INSERT INTO user_collection_item (user_id, collection_item_id) VALUES " + strings.Join(valueStrings, ",")

	_, err := tx.Exec(query, valueArgs...)
	return err
}

// レスポンス用のコレクション
func GetAllCollections(userID string) ([]*Collection, error) {
	userCollectionIDList, err := GetUserCollectionsIDList(userID)
	if err != nil {
		return nil, errors.Wrap(err, "[Error] failed to call GetUserCollectionsIDList")
	}

	collections, err := GetCollectionItems()
	if err != nil {
		return nil, errors.Wrap(err, "[Error] failed to call GetCollectionItems")
	}

	// TODO: 所持判定 レコード数が多い場合良くない（余裕があれば改善）
	for _, collection := range collections {
		if utils.Contains(userCollectionIDList, collection.CollectionID) {
			collection.HasItem = true
		}
	}

	return collections, nil
}
