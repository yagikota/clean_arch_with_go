package service

import "22dojo-online/pkg/domain/repository"

type SettingService interface {
}

type settingService struct {
	Repository repository.SettingRepository
}

func NewSettingService(settingRepository repository.SettingRepository) SettingService {
	return &settingService{
		Repository: settingRepository,
	}
}
