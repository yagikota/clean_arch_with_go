package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"

	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/db"
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"
	"22dojo-online/pkg/utils"
)

//nolint: gocyclo // this why
func HandleGachaDraw() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// log.Println("gacha")
		var requestBody gachadrawRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		times := requestBody.Times

		if times < 1 || times > constant.GachaTimesLimit {
			log.Println("times must be an positive interger, between 1 and 100")
			response.BadRequest(writer, "Bad Request")
			return
		}

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.BadRequest(writer, "userID is empty")
			return
		}

		// ユーザーが現在保有しているアイテムのIDリスト
		userCollectionIDList, err := model.GetUserCollectionsIDList(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// gacha_probabilityテーブルをもとに、出現率を考慮して出現するitemのIDを取得
		// ガチャで出現したアイテムのIDリスト
		gachaItemIDList, err := getGachaItemIDList(times)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// collection_itemテーブルから出現したitemをSELECT(重複なし)
		// QUESTION: getGachaItemIDLisに同じIDが複数含まれている場合の処理
		collections, err := model.GetCollectionItemsByID(gachaItemIDList)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// HasItem判定
		for _, collection := range collections {
			if utils.Contains(userCollectionIDList, collection.CollectionID) {
				collection.HasItem = true
			}
		}

		// 新しく獲得したアイテムだけを格納
		newCollectionIDList := make([]string, 0, times)
		for _, collection := range collections {
			if !collection.HasItem {
				newCollectionIDList = append(newCollectionIDList, collection.CollectionID)
			}
		}

		// QUESTION: transactionの範囲が適切か？
		// TODO: 発展課題以降で取り組む。
		// ==============================transaction start==============================
		err = db.ExecWithTx(ctx, db.Conn, func(tx *sql.Tx) error {
			user, err := model.SelectUserByPrimaryKeyWithLock(tx, userID)
			if err != nil {
				log.Println(err)
				return err
			}

			consumedCoin := constant.GachaCoinConsumption * times
			if consumedCoin > user.Coin {
				return errors.New("lack of coin")
			}

			// 出現したitemのうちHasItemがfalseのものだけ、user_collection_itemにINSERT
			if len(newCollectionIDList) > 0 {
				err = model.BulkInsertUserCollectionItemWithTx(tx, userID, newCollectionIDList)
				if err != nil {
					return err
				}
			}

			// userテーブルを更新（コイン減ったから）
			user.Coin -= consumedCoin
			err = model.UpdateUserByPrimaryKeyWithTx(tx, user)
			if err != nil {
				return err
			}
			return err
		})
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		// ==============================transaction end==============================

		// ガチャで出現したコレクションのIDとコレクションのmap（responseで使用）
		gachaCollectionMap := make(map[string]*model.Collection, len(collections))
		for _, collection := range collections {
			gachaCollectionMap[collection.CollectionID] = collection
		}

		results := make([]*result, 0, times)
		for _, id := range gachaItemIDList {
			collection, ok := gachaCollectionMap[id]
			if !ok {
				log.Println("gachaItemID not in gachaCollectionMap")
				response.InternalServerError(writer, "Internal Server Error")
				return
			}
			results = append(results, &result{
				CollectionID: collection.CollectionID,
				Name:         collection.Name,
				Rarity:       collection.Rarity,
				IsNew:        !collection.HasItem,
			})
		}

		response.Success(writer, &gachadrawResponse{
			Results: results,
		})
	}
}

// ガチャで出現するアイテムIDのリストを取得
// ratioを累積していき、生成した乱数(0 <= randInt < totalRatio)を初めて越えたアイテムを、出現アイテムと考える。
func getGachaItemIDList(times int32) ([]string, error) {
	gachaProbList, err := model.GetAllGachaProbs()
	if err != nil {
		return nil, err
	}

	var totalRatio int32
	for _, gachaProb := range gachaProbList {
		totalRatio += gachaProb.Ratio
	}

	getGachaItemIDList := make([]string, 0, times)
	for i := 0; i < int(times); i++ {
		//nolint: gosec // this is why
		randInt := rand.Int31n(totalRatio)

		var cumulativeRatio int32
		for _, gachaProb := range gachaProbList {
			cumulativeRatio += gachaProb.Ratio
			if randInt < cumulativeRatio {
				getGachaItemIDList = append(getGachaItemIDList, gachaProb.CollectionID)
				break
			}
		}
	}
	return getGachaItemIDList, nil
}

// response用のstruct
type gachadrawRequest struct {
	Times int32 `json:"times"`
}

type gachadrawResponse struct {
	Results []*result `json:"results"`
}

type result struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}
