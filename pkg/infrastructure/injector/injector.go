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
// User系======================================================
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

// User系======================================================

// Auth系======================================================
func InjectAuthController() middleware.AuthController {
	return middleware.NewAuthController(injectUserInteractor())
}

// Auth系======================================================

// Setting系===================================================
// TODO: これでいいのか
func injectSettingRepository() repository.SettingRepository {
	var tmpRepo repository.SettingRepository
	return tmpRepo
}

func injectSettingService() service.SettingService {
	return service.NewSettingService(injectSettingRepository())
}

func injectSettingInteractor() interactor.SettingrInteractor {
	return interactor.NewSettingInteractor(injectSettingService())
}

func InjectSettingController() controllers.SettingController {
	return controllers.NewSettingController(injectSettingInteractor())
}

// Collection系===================================================
func injectCollectionRepository() repository.CollectionRepository {
	return db.NewCollectionRepository(InjectDB())
}

func injectCollectionService() service.CollectionService {
	return service.NewCollectionService(injectCollectionRepository())
}

func injectCollectionInteractor() interactor.CollectionInteractor {
	return interactor.NewCollectionInteractor(injectCollectionService())
}

func InjectCollectionController() controllers.CollectionController {
	return controllers.NewCollectionController(injectCollectionInteractor())
}
