package core

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
)

// 原地修改meta的信息，当调用结束，meta应该是一个完整的目录
func ScanDir(cid string, meta *dir.Dir) {
	if meta == nil {
		meta = new(dir.Dir)
	}

	offset := 0
	for {
		dirInfo, err := ScanDirWithOffset(cid, offset)
		if err != nil {
			config.Logger.Warn(err.Error())
		}

		meta.DirName = dirInfo.Path[len(dirInfo.Path)-1].Name

		for _, item := range dirInfo.Data {
			if item.Fid != "" {
				// TODO 把任务通过channel派发出去

			}
		}
	}
}
