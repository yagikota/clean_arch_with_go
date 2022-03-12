package handler

import (
	"log"
	"net/http"

	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"

	"github.com/pkg/errors"
)

func HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// log.Println("collection/list API")
		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.BadRequest(writer, "User ID is empty")
			return
		}

		collections, err := model.GetAllCollections(userID)
		if err != nil {
			log.Printf("%+v\n", errors.Wrap(err, "[Error] failed to call GetAllCollections"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// メモリ確保
		collectionResponseList := make([]*collectionResponse, 0, len(collections))
		for _, collection := range collections {
			collectionResponseList = append(collectionResponseList, &collectionResponse{
				CollectionID: collection.CollectionID,
				Name:         collection.Name,
				Rarity:       collection.Rarity,
				HasItem:      collection.HasItem,
			})
		}
		response.Success(writer, &collectionsResponse{
			Collections: collectionResponseList,
		})
	}
}

// response用のstruct
type collectionsResponse struct {
	Collections []*collectionResponse `json:"collections"`
}

// response用のstruct
type collectionResponse struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}
