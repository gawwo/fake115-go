package executor

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
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

// 开启一定量的worker，通过channel接收任务，channel有一定的缓冲区
// worker在接收到任务后执行任务，当遇到需要人机验证的时候，改变全局
// 变量，然后进入循环等待模式，期间一直检测，直到人机验证完成；
// Note: 只要用到Lock的地方，都要考虑超时问题
// 将当前网络文件的内容导出到目录中
func (file *NetFile) Export(dir *dir.Dir) {
	downUrl := "http://webapi.115.com/files/download?pickcode=" + file.Pc
	headers := config.GetFakeHeaders(true)
	body, err := utils.Get(downUrl, headers, nil)
	if err != nil {
		config.Logger.Warn("export file network error",
			zap.String("name", file.Name))
		return
	}

	fmt.Println(body)
	return
}
