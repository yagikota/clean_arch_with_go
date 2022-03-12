package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"22dojo-online/pkg/db"
	"22dojo-online/pkg/dcontext"
	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"

	"github.com/pkg/errors"
)

func HandlerGameFinish() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody gameFinishRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Printf("%+v\n", errors.Wrap(err, "[Error] failed to decode requestbody"))
			response.BadRequest(writer, "Bad Request")
			return
		}

		score := requestBody.Score
		if score < 0 {
			log.Println("Score is negative integer")
			response.BadRequest(writer, "Bad Request")
			return
		}

		// log.Println("game finish API", score)

		if score == 0 {
			response.Success(writer, &gameFinishResponse{Coin: reward(score)})
			return
		}

		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.BadRequest(writer, "userID is empty")
			return
		}

		// TODO: トランザクション処理確認方法
		err := db.ExecWithTx(ctx, db.Conn, func(tx *sql.Tx) error {
			user, err := model.SelectUserByPrimaryKeyWithLock(tx, userID)
			if err != nil {
				return errors.Wrap(err, "[Error] failed to call SelectUserByPrimaryKeyWithLock")
			}

			// ハイスコア更新
			if user.HighScore < score {
				user.HighScore = score
			}
			// コイン枚数更新
			user.Coin += reward(score)

			return model.UpdateUserByPrimaryKeyWithTx(tx, user)
		})
		if err != nil {
			log.Printf("%+v\n", errors.Wrap(err, "[Error] failed to call ExecWithTx"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		response.Success(writer, &gameFinishResponse{Coin: reward(score)})
	}
}

// 報酬のコインの計算式
func reward(score int32) int32 {
	return score
}

type gameFinishRequest struct {
	Score int32 `json:"score"`
}

type gameFinishResponse struct {
	Coin int32 `json:"coin"`
}
