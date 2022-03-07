package handler

import (
	"log"
	"net/http"
	"strconv"

	"22dojo-online/pkg/http/response"
	"22dojo-online/pkg/server/model"

	"github.com/pkg/errors"
)

func HandleRankingList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// queryからstart取得
		start, err := strconv.Atoi(request.URL.Query().Get("start"))
		if err != nil {
			log.Printf("%+v\n", errors.Wrap(err, "[Error] failed to convert string to int"))
			response.BadRequest(writer, "Bad Request, start query is invalid")
			return
		}
		// log.Println("ranking/list API", start)
		
		if start <= 0 {
			log.Println(errors.New("start must be an integer, between 1 and the number of user"))
			response.BadRequest(writer, "Bad Request: start must be an integer, between 1 and the number of user")
			return
		}

		rankingList, err := model.GetRanking(start)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		ranks := make([]*rankResponse, 0, len(rankingList))
		for _, rank := range rankingList {
			ranks = append(ranks, &rankResponse{
				UserID:   rank.UserID,
				UserName: rank.UserName,
				Rank:     rank.Rank,
				Score:    rank.Score,
			})
		}
		response.Success(writer, &rankingListResponse{
			Ranks: ranks,
		})
	}
}

type rankingListResponse struct {
	Ranks []*rankResponse `json:"ranks"`
}

type rankResponse struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}
