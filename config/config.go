package config

import "fake115/log"

var (
	ServerName = "fake115"
	AppVer     = "11.2.0"
	EndString  = "000000"
	UserAgent  = "Mozilla/5.0  115disk/11.2.0"
	RetryTimes = 3
	UserId = ""
	UserKey = ""

	// Cookie提供文件读取和命令行设置两种方式；
	// 文件读取提供默认设置和命令行设置两种方式；
	Cookie     = ""
	DefaultCookiePath = "cookies.txt"
	CookiePath = ""

	Logger = log.InitLogger(ServerName, false)
)

var fakeHeaders = map[string]string{
	"User-Agent": UserAgent,
}

func GetFakeHeaders(withCookie bool) map[string]string {

	copyMap := map[string]string{}
	for k, v := range fakeHeaders {
		copyMap[k] = v
	}

	if withCookie {
		copyMap["Cookie"] = Cookie
	}

	return copyMap
}

var fakeRangeHeaders = map[string]string{
	"Range":      "bytes=0-131071",
	"User-Agent": UserAgent,
}

func GetFakeRangeHeaders() map[string]string {
	copyMap := map[string]string{}
	for k, v := range fakeRangeHeaders {
		copyMap[k] = v
	}
	return copyMap
}
