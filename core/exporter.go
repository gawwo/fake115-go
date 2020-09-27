package core

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
	"time"
)

// 原地修改meta的信息，当调用结束，meta应该是一个完整的目录
func scanDir(cid string, meta *dir.Dir, sem *utils.WaitGroupPool) {
	// 递归调用的初始调用，与其他这个函数的递归调用不一样；
	// 初始的调用，需要等待一会，让其他递归的这个函数拿到信
	// 号量，递归的则需要放回自己的信号量；
	var newest = false

	defer func() {
		if !newest {
			sem.Done()
		} else {
			// 给迭代的scanDir一点获取信号量的时间
			time.Sleep(time.Second)
		}
	}()

	if sem == nil {
		sem = dir.WaitGroupPool
		newest = true
	} else {
		// 太多的scanDir worker会导致阻塞，用以避免scanDir数量失控
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
				innerMeta := dir.NewDir()
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
	meta := dir.NewDir()
	scanDir(cid, meta, nil)

	// 等待生产者资源枯竭之后，关闭channel
	dir.WaitGroupPool.Wait()
	close(WorkerChannel)

	// 等待消费者完成任务
	config.WaitGroup.Wait()

	return meta
}
