package core

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
)

// 原地修改meta的信息，当调用结束，meta应该是一个完整的目录
func scanDir(cid string, meta *dir.Dir, sem *utils.WaitGroupPool) {
	defer func() {
		if sem != nil {
			sem.Done()
		}
	}()

	if meta == nil {
		meta = new(dir.Dir)
	}
	if sem == nil {
		sem = dir.WaitGroupPool
	} else {
		// 太多的scanDir worker会导致阻塞，以避免scanDir数量失控
		sem.Add()
	}

	offset := 0
	for {
		dirInfo, err := ScanDirWithOffset(cid, offset)
		if err != nil {
			config.Logger.Warn(err.Error())
			return
		}

		meta.DirName = dirInfo.Path[len(dirInfo.Path)-1].Name

		for _, item := range dirInfo.Data {
			if item.Fid != "" {
				// 处理文件
				// 把任务通过channel派发出去
				task := Task{Dir: meta, File: item}
				WorkerChannel <- task
			} else if item.Cid != "" {
				// 处理文件夹
				innerMeta := new(dir.Dir)
				meta.Dirs = append(meta.Dirs, innerMeta)
				go scanDir(item.Cid, innerMeta, sem)
			}
		}

		// 翻页
		if dirInfo.Count-(offset+config.Step) > 0 {
			offset += config.Step
			continue
		} else {
			config.Logger.Info("scan dir success", zap.String("name", meta.DirName))
			return
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
	// meta是提取资源的抓手
	meta := new(dir.Dir)
	scanDir(cid, meta, nil)

	// 生产者资源枯竭
	close(WorkerChannel)

	// 等待消费者完成任务
	config.WaitGroup.Wait()

	return meta
}
