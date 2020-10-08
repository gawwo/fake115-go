package core

import (
	"encoding/json"
	"errors"
	"fmt"
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

func ScanDirWithOffset(cid string, offset int) (*NetDir, error) {
	urls := []string{
		fmt.Sprintf("https://webapi.115.com/files?aid=1&cid=%s&o=file_name&asc=0&offset=%d&"+
			"show_dir=1&limit=%d&code=&scid=&snap=0&natsort=1&record_open_time=1&source=&format=json&type=&star=&"+
			"is_share=&suffix=&fc_mix=1&is_q=&custom_order=", cid, offset, config.Step),
		fmt.Sprintf("http://aps.115.com/natsort/files.php?aid=1&cid=%s&o=file_name&asc=1&offset=%d&show_dir=1"+
			"&limit=%d&code=&scid=&snap=0&natsort=1&source=&format=json&type=&star=&is_share=&suffix=&custom_order="+
			"&fc_mix=", cid, offset, config.Step),
	}

	for _, url := range urls {
		headers := config.GetFakeHeaders(true)
		body, err := utils.Get(url, headers, nil)
		if err != nil {
			continue
		}

		netDir := new(NetDir)
		err = json.Unmarshal(body, netDir)
		if err != nil {
			continue
		}
		if netDir.Path == nil {
			continue
		} else {
			return netDir, nil
		}
	}
	return nil, errors.New("both url fail get dir info, maybe login expire")
}
