package interactor

import (
	"22dojo-online/pkg/domain/service"
)

type SettingrInteractor interface {
}

type settingInteractor struct {
	Service service.SettingService
}

func NewSettingInteractor(settingService service.SettingService) SettingrInteractor {
	return &settingInteractor{
		Service: settingService,
	}
}
