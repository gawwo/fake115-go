package core

import (
	"fmt"
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
		if config.Debug {
			fmt.Println("Dir digger on work number: ", sem.Size())
		}
		// 防止goroutine过早的退出，过早退出会导致sem的Wait可能过早的
		// 返回，但实际上下一个goroutine还没有Add到信号量，Wait
		// 返回后还会导致传递task的通道关闭，进而导致整个任务提早结束
		if sem.Size() <= 1 {
			time.Sleep(time.Second * 2)
		}

		if !newest {
			sem.Done()
		}
	}()

	if sem == nil {
		sem = dir.ProducerWaitGroupPool
		newest = true
	} else {
		// 太多的scanDir worker会导致阻塞，用以避免scanDir数量失控
		sem.Add()
	}

	defer func() {
		if config.Debug {
			fmt.Println("Dir digger on work number: ", sem.Size())
		}
		// 防止goroutine过早的退出，过早退出会导致sem的Wait可能过早的
		// 返回，但实际上下一个goroutine还没有Add到信号量，Wait
		// 返回后还会导致传递task的通道关闭，进而导致整个任务提早结束
		if sem.Size() <= 1 {
			time.Sleep(time.Second * 2)
		}

		if !newest {
			sem.Done()
		}
	}()

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
				task := ExportTask{Dir: meta, File: item}
				ExportWorkerChannel <- task
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
	config.ConsumerWaitGroup.Add(config.WorkerNum)
	for i := 0; i < config.WorkerNum; i++ {
		go ExportWorker()
	}

	// 开启生产者
	// meta是提取资源的抓手
	meta := dir.NewDir()
	scanDir(cid, meta, nil)

	// 等待生产者资源枯竭之后，关闭channel
	dir.ProducerWaitGroupPool.Wait()
	close(ExportWorkerChannel)

	// 等待消费者完成任务
	config.ConsumerWaitGroup.Wait()

	return meta
}

func Export(cid string) {
	dirMeta := ScanDir(cid)
	exportName := fmt.Sprintf("%s_%s_%dGB.json", cid, dirMeta.DirName,
		config.TotalSize>>30)
	_, err := dirMeta.Dump(exportName)
	if err != nil {
		fmt.Println("导出到文件失败")
		return
	}

	fmt.Printf("导出文件%s成功, 文件数： %d\n", exportName, config.FileCount)
}
