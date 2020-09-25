package executor

import (
	"encoding/json"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
)

type userInfo struct {
	UserId  int    `json:"user_id"`
	UserKey string `json:"userkey"`
	ErrorNo int    `json:"errno"`
}

func SetUserInfoConfig() bool {
	url := "http://proapi.115.com/app/uploadinfo"
	header := config.GetFakeHeaders(true)
	body, err := utils.Get(url, header, nil)
	if err != nil {
		return false
	}

	jsonUserInfo := new(userInfo)
	err = json.Unmarshal(body, jsonUserInfo)
	if err != nil {
		config.Logger.Error("Get login information error",
			zap.Binary("content", body))
		return false
	}

	if jsonUserInfo.ErrorNo == 99 {
		config.Logger.Error("Login expire")
		return false
	} else if jsonUserInfo.UserKey == "" {
		config.Logger.Error("Get login information error",
			zap.String("content", string(body)))
		return false
	}

	// 设置用户信息
	config.UserId = jsonUserInfo.UserId
	config.UserKey = jsonUserInfo.UserKey

	return true
}
