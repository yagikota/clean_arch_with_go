package controllers

import (
	"log"
	"net/http"

	"22dojo-online/pkg/adapter/dcontext"
	"22dojo-online/pkg/adapter/response"
	interactor "22dojo-online/pkg/usecase/Interactor"
	outputdata "22dojo-online/pkg/usecase/output_data"

	"github.com/pkg/errors"
)

type CollectionController interface {
	HandleCollectionList() http.HandlerFunc
}

type collectionController struct {
	Interactor interactor.CollectionInteractor
}

func NewCollectionController(collectionInteractor interactor.CollectionInteractor) CollectionController {
	return &collectionController{
		Interactor: collectionInteractor,
	}
}

func (cc *collectionController) HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.BadRequest(writer, "User ID is empty")
			return
		}

		collections, err := cc.Interactor.GetAllCollections(userID)
		if err != nil {
			log.Printf("%+v\n", errors.Wrap(err, "[Error] failed to call GetAllCollections"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// メモリ確保
		collectionResponseList := make([]*outputdata.CollectionResponse, 0, len(collections))
		for _, collection := range collections {
			collectionResponseList = append(collectionResponseList, &outputdata.CollectionResponse{
				CollectionID: collection.CollectionID,
				Name:         collection.Name,
				Rarity:       collection.Rarity,
				HasItem:      collection.HasItem,
			})
		}

		response.Success(writer, &outputdata.CollectionsResponse{
			Collections: collectionResponseList,
		})
	}
}
