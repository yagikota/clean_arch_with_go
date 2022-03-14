package controllers

import (
	"net/http"

	"22dojo-online/pkg/adapter/response"
	"22dojo-online/pkg/constant"
	interactor "22dojo-online/pkg/usecase/Interactor"
	outputdata "22dojo-online/pkg/usecase/output_data"
)

type SettingController interface {
	HandleSettingGet() http.HandlerFunc
}

type settingController struct {
	Interactor interactor.SettingrInteractor
}

func NewSettingController(settingInteractor interactor.SettingrInteractor) SettingController {
	return &settingController{
		Interactor: settingInteractor,
	}
}

// 汎用的な実装になるように他のAPIと同様の実装にする。
func (controller *settingController) HandleSettingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response.Success(writer, &outputdata.SettingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
		})
	}
}
