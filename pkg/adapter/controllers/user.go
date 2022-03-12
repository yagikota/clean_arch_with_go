package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"22dojo-online/pkg/adapter/dcontext"
	"22dojo-online/pkg/adapter/response"
	interactor "22dojo-online/pkg/usecase/Interactor"
	inputdata "22dojo-online/pkg/usecase/input_data"
	outputdata "22dojo-online/pkg/usecase/output_data"
)

type UserController struct {
	Interactor interactor.UserInteractor
}

func NewUserController(userInteractor interactor.UserInteractor) *UserController {
	return &UserController{
		Interactor: userInteractor,
	}
}

func (controller *UserController) HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// TODO: interfaceにしたい
		var requestBody inputdata.UserCreateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		authToken, err := controller.Interactor.CreateUser(requestBody)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		response.Success(writer, &outputdata.UserCreateResponse{Token: authToken})
	}
}

func (controller *UserController) HandleUserGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("user get")
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.BadRequest(writer, "userID is empty")
			return
		}
		user, err := controller.Interactor.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println("user not found")
			response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}
		response.Success(writer, &outputdata.UserGetResponse{
			ID:        user.ID,
			Name:      user.Name,
			HighScore: user.HighScore,
			Coin:      user.Coin,
		})
	}
}

func (controller *UserController) HandleUserUpdate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody inputdata.UserUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
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

		if err := controller.Interactor.UpdateUserByPrimaryKey(requestBody, userID); err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		response.Success(writer, nil)
	}
}
