package interactor

import (
	"22dojo-online/pkg/domain/model"
	"22dojo-online/pkg/domain/service"
)

// inputdataをもとに処理

type CollectionInteractor interface {
	GetUserCollectionsIDList(userID string) ([]string, error)
	GetCollectionItems() ([]*model.Collection, error)
	GetCollectionItemsByID(idList []string) ([]*model.Collection, error)
	// BulkCreateUserCollectionItemWithTx(tx *sql.Tx, userID string, idList []string) error
	GetAllCollections(userID string) ([]*model.Collection, error)
}

type collectionInteractor struct {
	Service service.CollectionService
}

func NewCollectionInteractor(collectionService service.CollectionService) CollectionInteractor {
	return &collectionInteractor{
		Service: collectionService,
	}
}

func (ci *collectionInteractor) GetUserCollectionsIDList(userID string) ([]string, error) {
	return ci.Service.GetUserCollectionsIDList(userID)
}

func (ci *collectionInteractor) GetCollectionItems() ([]*model.Collection, error) {
	return ci.Service.GetCollectionItems()
}

func (ci *collectionInteractor) GetCollectionItemsByID(idList []string) ([]*model.Collection, error) {
	return ci.Service.GetCollectionItemsByID(idList)
}

func (ci *collectionInteractor) GetAllCollections(userID string) ([]*model.Collection, error) {
	return ci.Service.GetAllCollections(userID)
}
