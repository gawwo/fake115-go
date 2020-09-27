package core

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
)

// 原地修改meta的信息，当调用结束，meta应该是一个完整的目录
func scanDir(cid string, meta *dir.Dir) {
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
				// 把任务通过channel派发出去
				task := Task{Dir: meta, File: item}
				config.WorkerChannel <- task
			}
		}
	}
}

func ScanDir(cid string) *dir.Dir {
	// 开启消费者
	config.WaitGroup.Add(config.WorkerNum)
	for i := 0; i < config.WorkerNum; i++ {
		go Worker()
	}

	// 开启生产者
	// meta是提取资源的句柄
	meta := new(dir.Dir)
	scanDir(cid, meta)

	// 生产者资源枯竭
	close(config.WorkerChannel)

	// 等待消费者完成任务
	config.WaitGroup.Wait()

	return meta
}
