package service

import (
	"22dojo-online/pkg/domain/model"
	"22dojo-online/pkg/domain/repository"
)

type CollectionService interface {
	GetUserCollectionsIDList(userID string) ([]string, error)
	GetCollectionItems() ([]*model.Collection, error)
	GetCollectionItemsByID(idList []string) ([]*model.Collection, error)
	// BulkCreateUserCollectionItemWithTx(tx *sql.Tx, userID string, idList []string) error
	GetAllCollections(userID string) ([]*model.Collection, error)
}

type collectionService struct {
	Repository repository.CollectionRepository
}

func NewCollectionService(collectionRepository repository.CollectionRepository) CollectionService {
	return &collectionService{
		Repository: collectionRepository,
	}
}

func (cs *collectionService) GetUserCollectionsIDList(userID string) ([]string, error) {
	return cs.Repository.GetUserCollectionsIDList(userID)
}

func (cs *collectionService) GetCollectionItems() ([]*model.Collection, error) {
	return cs.Repository.GetCollectionItems()
}

func (cs *collectionService) GetCollectionItemsByID(idList []string) ([]*model.Collection, error) {
	return cs.Repository.GetCollectionItemsByID(idList)
}

func (cs *collectionService) GetAllCollections(userID string) ([]*model.Collection, error) {
	return cs.Repository.GetAllCollections(userID)
}
