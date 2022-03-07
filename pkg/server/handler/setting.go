package handler

import (
	"net/http"

	"22dojo-online/pkg/constant"
	"22dojo-online/pkg/http/response"
)

// HandleSettingGet ゲーム設定情報取得処理
func HandleSettingGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response.Success(writer, &settingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
		})
	}
}

type settingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
}
