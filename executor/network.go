package executor

import (
	"encoding/json"
	"fake115/config"
	"fake115/utils"
)

type userInfo struct {
	userId  string `json:"user_id"`
	userKey string `json:"userkey"`
}

func getUserInfo() bool {
	url := "http://proapi.115.com/app/uploadinfo"
	header := config.GetFakeHeaders(true)
	body, err := utils.Get(url, header, nil)
	if err != nil {
		return false
	}

	jsonUserInfo := new(userInfo)
	err = json.Unmarshal(body, jsonUserInfo)
	if err != nil {
		return false
	}
	return true
}
