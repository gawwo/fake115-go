package utils

import (
	"github.com/gawwo/fake115-go/config"
	"io/ioutil"
	"os"
)

// ReadCookieFile 一定要在main中的init函数之后运行
func ReadCookieFile() (string, error) {
	cookieFile, err := os.Open(config.CookiePath)
	if err != nil {
		if os.IsNotExist(err) {
			config.Logger.Warn("Fail to open cookie dir: " + err.Error())
		}
		return "", err
	}
	defer cookieFile.Close()

	cookie, err := ioutil.ReadAll(cookieFile)
	if err != nil {
		return "", err
	}
	return string(cookie), nil
}
