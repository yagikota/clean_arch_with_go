package server

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"22dojo-online/pkg/adapter/controllers"
	"22dojo-online/pkg/adapter/db"
	"22dojo-online/pkg/adapter/middleware"
	infrasql "22dojo-online/pkg/infrastructure/sql"
	interactor "22dojo-online/pkg/usecase/Interactor"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	rand.Seed(time.Now().UnixNano())
	// 依存性の注入(Constructor Injection)
	sqlHandler := infrasql.NewSQLHandler()
	// userRepo := db.NewUserRepository(sqlHandler)
	userRepo := db.NewUserRepository(sqlHandler)
	userInteractor := interactor.NewUserInteractor(userRepo)
	userController := controllers.NewUserController(userInteractor)
	
	// userController := controllers.NewUserController(&db.NewSQLHandler())
	// http.HandleFunc("/setting/get", get(controllers.HandleSettingGet()))
	http.HandleFunc("/user/create", post(userController.HandleUserCreate()))
	http.HandleFunc("/user/get",
		get(middleware.Authenticate(userController.HandleUserGet())))
	http.HandleFunc("/user/update",
		post(middleware.Authenticate(userController.HandleUserUpdate())))
	// http.HandleFunc("/game/finish",
	// 	post(middleware.Authenticate(controllers.HandlerGameFinish())))

	// http.HandleFunc("/ranking/list",
	// 	get(middleware.Authenticate(controllers.HandleRankingList())))
	// http.HandleFunc("/collection/list",
	// 	get(middleware.Authenticate(controllers.HandleCollectionList())))
	// http.HandleFunc("/gacha/draw",
	// 	post(middleware.Authenticate(controllers.HandleGachaDraw())))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			if _, err := writer.Write([]byte("Method Not Allowed")); err != nil {
				log.Println(err)
			}
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}

func InjectDB() *db.DBHandler {
	return db.NewDBHandler()
}

func injectUserRepository() repository.UserRepository {
	return mysqlhandler.NewUserRepository(InjectDB())
}

func injectUserService() service.UserService {
	return service.NewUserService(injectUserRepository())
}

func injectUserInteractor() interactor.UserInteractor {
	return interactor.NewUserInteractor(injectUserService())
}

func injectUserHandler() usecasehandler.UserHandler {
	return usecasehandler.NewUserHandler(injectUserInteractor())
}
