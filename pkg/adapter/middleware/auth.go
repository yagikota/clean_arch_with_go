package middleware

import (
	"context"
	"log"
	"net/http"

	"22dojo-online/pkg/adapter/dcontext"
	"22dojo-online/pkg/adapter/response"
	interactor "22dojo-online/pkg/usecase/Interactor"
)

type AuthController struct {
	Interactor interactor.UserInteractor
}

func NewAuthController(userInteractor interactor.UserInteractor) *AuthController {
	return &AuthController{
		Interactor: userInteractor,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (auth *AuthController) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			log.Println("x-token is empty")
			return
		}

		// TODO: データベースから認証トークンに紐づくユーザの情報を取得
		user, err := auth.Interactor.SelectUserByAuthToken(token)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Invalid token")
			return
		}
		if user == nil {
			log.Printf("user not found. token=%s", token)
			response.BadRequest(writer, "Invalid token")
			return
		}

		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 次の処理
		nextFunc(writer, request.WithContext(ctx))
	}
}
