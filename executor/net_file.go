package executor

import (
	"encoding/json"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
	"time"
)

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
func (file *NetFile) Export(dir *dir.Dir) {
	// 保证worker不会panic
	defer func() {
		if err := recover(); err != nil {
			config.Logger.Error("export link error",
				zap.String("reason", fmt.Sprintf("%v", err)))

			// 在报错的情况下，如果依然处于人机验证的阻塞状态，就解除状态
			if config.SpiderVerification {
				config.SpiderVerification = false
			}
		}
	}()

	downUrl := "http://webapi.115.com/files/download?pickcode=" + file.Pc
	headers := config.GetFakeHeaders(true)
	for {
		// 先检查是否在等待人机验证状态
		headOff := config.SpiderVerification
		if headOff {
			config.Logger.Info(fmt.Sprintf("waiting Man-machine verification: %s", file.Name))
			time.Sleep(time.Second * 10)
			continue
		}

	Work:
		body, err := utils.Get(downUrl, headers, nil)
		if err != nil {
			config.Logger.Warn("export file network error",
				zap.String("name", file.Name))
			return
		}

		fmt.Println(body)
		parsedDownloadBody := new(downloadBody)
		err = json.Unmarshal(body, parsedDownloadBody)
		if err != nil {
			config.Logger.Warn("parse download body fail",
				zap.String("content", string(body)))
			return
		}

		// 文件状态异常
		if !parsedDownloadBody.State {
			config.Logger.Warn("download file state odd",
				zap.String("content", parsedDownloadBody.Msg))
			return
		}

		// TODO 人机验证处理xxx改成触发值
		// 有多个worker因为时间差，都进入人机检测验证状态，也无所谓
		if parsedDownloadBody.Msg == "xxx" {
			config.Logger.Warn("found Man-machine verification， waiting...")
			config.SpiderVerification = true
			time.Sleep(time.Second * 15)
			goto Work
		}

		// 返回的下载信息中不包含下载地址
		if parsedDownloadBody.FileUrl == "" {
			config.Logger.Warn("download file body not contain download url",
				zap.String("content", fmt.Sprintf("%v", parsedDownloadBody)))
		}

		// 如果没有人机验证,就尝试取消人机验证状态
		if config.SpiderVerification {
			config.SpiderVerification = false
		}

		return
	}

}
