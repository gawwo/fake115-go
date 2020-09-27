package core

import (
	"encoding/json"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock sync.Mutex

// 115的文件对象，这个对象指向的可能是文件，也可能是文件夹
type NetFile struct {
	// 有fid就是文件
	Fid string `json:"fid"`
	// 有cid但没有Fid，就是文件夹
	Cid string `json:"cid"`
	// 文件大小
	Size int    `json:"s"`
	Name string `json:"n"`
	Sha  string `json:"sha"`
	Pc   string `json:"pc"`
}

type downloadBody struct {
	State   bool   `json:"state"`
	Msg     string `json:"msg"`
	FileUrl string `json:"file_url"`
}

// 开启一定量的worker，通过channel接收任务，channel有一定的缓冲区
// worker在接收到任务后执行任务，当遇到需要人机验证的时候，改变全局
// 变量，然后进入循环等待模式，期间一直检测，直到人机验证完成；
// Note: 只要用到Lock的地方，都要考虑超时问题
// 将当前网络文件的内容导出到目录中
func (file *NetFile) Export() string {
	// 保证worker不会panic
	defer func() {
		if err := recover(); err != nil {
			config.Logger.Error("export link error",
				zap.String("reason", fmt.Sprintf("%v", err)))

			// 在报错的情况下，如果依然处于人机验证的阻塞状态，就解除状态。不
			// 希望因为一个处于检测状态的任务panic而导致整个任务卡死
			if config.SpiderVerification {
				config.SpiderVerification = false
			}
		}
	}()

	url, cookie := file.extractDownloadInfo()
	if cookie == "" || url == "" {
		return ""
	}

	fileSha1 := file.extractFileSha1(url, cookie)
	if fileSha1 == "" {
		return ""
	}

	joinStrings := []string{file.Name, strconv.Itoa(file.Size), file.Sha, fileSha1}
	result := strings.Join(joinStrings, "|")

	return result
}

func (file *NetFile) extractDownloadInfo() (downloadUrl, cookie string) {
	downUrl := "http://webapi.115.com/files/download?pickcode=" + file.Pc
	headers := config.GetFakeHeaders(true)
	for {
		// 先检查是否在等待人机验证状态
		headOff := config.SpiderVerification
		if headOff {
			config.Logger.Info(fmt.Sprintf("waiting Man-machine verification: %s", file.Name))
			time.Sleep(config.SpiderCheckInterval / 2)
			continue
		}

	Work:
		body, response, err := utils.GetResponse(downUrl, headers, nil)
		if err != nil {
			config.Logger.Warn("export file network error",
				zap.String("name", file.Name))
			return
		}

		parsedDownloadBody := new(downloadBody)
		err = json.Unmarshal(body, parsedDownloadBody)
		if err != nil {
			config.Logger.Warn("parse download body fail",
				zap.String("content", string(body)),
				zap.String("name", file.Name))
			return
		}

		// ==================验证文件下载地址是否正常==================
		// 文件状态异常
		if !parsedDownloadBody.State {
			config.Logger.Warn("download file state odd",
				zap.String("content", parsedDownloadBody.Msg),
				zap.String("name", file.Name))
			return
		}

		// TODO 人机验证处理xxx改成触发值
		// 有多个worker因为时间差，都进入人机检测验证状态，也无所谓
		// 进入人机验证之后，反复检测状态
		if parsedDownloadBody.Msg == "xxx" {
			fmt.Println("发现人机验证，请到115浏览器中播放任意一个视频，完成人机检测...")
			config.Logger.Warn("found Man-machine verification， waiting...")
			config.SpiderVerification = true
			time.Sleep(config.SpiderCheckInterval)
			goto Work
		}
		// 如果没有人机验证,就尝试取消人机验证状态
		if config.SpiderVerification {
			config.SpiderVerification = false
		}

		// 返回的下载信息中不包含下载地址
		if parsedDownloadBody.FileUrl == "" {
			config.Logger.Warn("download file body not contain download url",
				zap.String("content", fmt.Sprintf("%v", parsedDownloadBody)),
				zap.String("name", file.Name))
			return
		}
		// 下载的时候有自己单独的cookie,提取下载cookie
		cookie := downloadCookie(response)
		if cookie == "" {
			config.Logger.Warn("get download cookie fail", zap.String("name", file.Name))
			return "", ""
		}
		return parsedDownloadBody.FileUrl, cookie
	}
}

func downloadCookie(response *http.Response) string {
	newCookie, ok := response.Header["Set-Cookie"]
	if ok && len(newCookie) >= 1 {
		cookies := strings.SplitN(newCookie[0], ";", 2)
		if len(cookies) >= 2 {
			return cookies[0]
		}
	}
	return ""
}

func (file *NetFile) extractFileSha1(downloadUrl, cookie string) string {
	downloadHeader := config.GetFakeRangeHeaders()
	downloadHeader["Cookie"] = cookie

	body, err := utils.Get(downloadUrl, downloadHeader, nil)
	if err != nil {
		config.Logger.Warn("get file header to calculate file sha1 fail", zap.String("name", file.Name))
		return ""
	}

	sha1 := utils.Sha1(body)
	return strings.ToUpper(sha1)
}
