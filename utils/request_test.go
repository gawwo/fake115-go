package utils

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	config.CookiePath = config.DefaultCookiePath
	config.Cookie, _ = ReadCookieFile()

	header := config.GetFakeHeaders(true)
	url := "https://webapi.115.com/files"
	data := map[string]string{
		"aid":              "1",
		"cid":              "353522044329243945",
		"o":                "file_name",
		"asc":              "0",
		"offset":           "0",
		"show_dir":         "1",
		"limit":            "115",
		"code":             "",
		"scid":             "",
		"snap":             "0",
		"natsrot":          "1",
		"record_open_time": "1",
		"source":           "",
		"format":           "json",
		"type":             "",
		"star":             "",
		"is_share":         "",
		"suffix":           "",
		"is_q":             "",
		"fc_mix":           "1",
	}

	body, err := Get(url, header, data)
	if err != nil {
		t.Errorf("请求错误： %s\n", err)
	} else {
		fmt.Println(body)
	}
}

var wgp = NewWaitGroupPool(50)

func doMultiRequest(url string) bool {
	defer wgp.Done()

	start := time.Now().Unix()
	_, err := Get(url, nil, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("elapsed: ", time.Now().Unix()-start)
	return true
}

func TestTimeout(t *testing.T) {
	url := "http://localhost:8441/"
	for i := 0; i < 5000000; i++ {
		wgp.Add()
		go doMultiRequest(url)
	}

	wgp.Wait()
}

func TestPostForm(t *testing.T) {
	config.CookiePath = config.DefaultCookiePath
	config.Cookie, _ = ReadCookieFile()

	header := config.GetFakeHeaders(true)
	url := "https://webapi.115.com/files/add"
	data := map[string]string{
		"pid":   "353522044329243945",
		"cname": "1",
	}

	body, err := PostForm(url, header, data)
	if err != nil {
		t.Errorf("请求错误： %s\n", err)
	} else {
		fmt.Println(body)
	}
}
