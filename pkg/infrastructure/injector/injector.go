package injector

import (
	"22dojo-online/pkg/adapter/controllers"
	"22dojo-online/pkg/adapter/db"
	"22dojo-online/pkg/adapter/middleware"
	"22dojo-online/pkg/domain/repository"
	"22dojo-online/pkg/domain/service"
	infrasql "22dojo-online/pkg/infrastructure/sql"
	interactor "22dojo-online/pkg/usecase/Interactor"
)

// TODO: go wireなどで自動生成
// 依存するオブジェクトを注入
func InjectDB() *infrasql.SQLHandler {
	return infrasql.NewSQLHandler()
}

// TODO: 実装方法?
func injectUserRepository() repository.UserRepository {
	return db.NewUserRepository(InjectDB())
}

func injectUserService() service.UserService {
	return service.NewUserService(injectUserRepository())
}

func injectUserInteractor() interactor.UserInteractor {
	return interactor.NewUserInteractor(injectUserService())
}

func InjectUserController() controllers.UserController {
	return controllers.NewUserController(injectUserInteractor())
}

func InjectAuthController() middleware.AuthController {
	return middleware.NewAuthController(injectUserInteractor())
}
