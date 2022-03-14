package repository

import (
	"22dojo-online/pkg/domain/model"
)

// TODO: sql処理ないで呼ばれているsql処理についての扱い
type CollectionRepository interface {
	GetUserCollectionsIDList(userID string) ([]string, error)
	GetCollectionItems() ([]*model.Collection, error)
	GetCollectionItemsByID(idList []string) ([]*model.Collection, error)
	// BulkCreateUserCollectionItemWithTx(tx *sql.Tx, userID string, idList []string) error
	GetAllCollections(userID string) ([]*model.Collection, error)
}
