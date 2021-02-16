package core

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"time"
)

type ExportTask struct {
	Dir  *dir.Dir
	File *NetFile
}

type exporter struct {
	taskChannel       chan ExportTask
	consumerWaitGroup sync.WaitGroup
	// 通过pool支持设置上限
	producerWaitGroupPool *utils.WaitGroupPool
	lock                  sync.Mutex
	FileCount             int
	FileTotalSize         int
}

func NewExporter() *exporter {
	return &exporter{
		taskChannel:           make(chan ExportTask, config.WorkerNum*config.WorkerNumRate),
		consumerWaitGroup:     sync.WaitGroup{},
		producerWaitGroupPool: utils.NewWaitGroupPool(config.WorkerNum),
		FileCount:             0,
		FileTotalSize:         0,
	}
}

// 原地修改meta的信息，当调用结束，meta应该是一个完整的目录
func (e *exporter) scanDir(cid string, meta *dir.Dir) {
	defer func() {
		if config.Debug {
			fmt.Println("Dir digger on work number: ", e.producerWaitGroupPool.Size())
		}
		// 防止goroutine过早的退出，过早退出会导致sem的Wait可能过早的
		// 返回，但实际上下一个goroutine还没有Add到信号量，Wait
		// 返回后还会导致传递task的通道关闭，进而导致整个任务提早结束
		runtime.Gosched()
		e.producerWaitGroupPool.Done()
	}()

	e.producerWaitGroupPool.Add()

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
				e.taskChannel <- task
			} else if item.Cid != "" {
				// 处理文件夹
				innerMeta := dir.NewDir()
				meta.Dirs = append(meta.Dirs, innerMeta)
				go e.scanDir(item.Cid, innerMeta)
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

func (e *exporter) ScanDir(cid string) *dir.Dir {
	// 开启消费者
	e.consumerWaitGroup.Add(config.WorkerNum)
	for i := 0; i < config.WorkerNum; i++ {
		go e.exportConsumer()
	}

	// 开启生产者
	// meta是提取资源的抓手
	meta := dir.NewDir()
	e.scanDir(cid, meta)

	// 等待生产者资源枯竭之后，关闭channel
	e.producerWaitGroupPool.Wait()
	close(e.taskChannel)

	// 等待消费者完成任务
	e.consumerWaitGroup.Wait()

	return meta
}

func (e *exporter) exportConsumer() {
	// WorkerChannel关闭前一直工作，直到生产者枯竭
	for task := range e.taskChannel {
		if config.Debug {
			fmt.Println("channel len: ", len(e.taskChannel))
		}
		// 有recover，保证这里不会panic，能让任务持续进行
		start := time.Now().Unix()
		result := task.File.Export()
		if result == "" {
			config.Logger.Warn("export failed", zap.String("name", task.File.Name))
			continue
		}

		// 监控时间太长的请求
		elapsed := time.Now().Unix() - start
		if elapsed > int64(3) {
			config.Logger.Warn("task slow", zap.String("name", task.File.Name),
				zap.Int64("elapsed", elapsed))
		}

		// 扫尾工作，添加记录到dir对象，累加文件总大小
		e.lock.Lock()
		task.Dir.Files = append(task.Dir.Files, result)
		e.FileTotalSize += task.File.Size
		e.FileCount += 1
		e.lock.Unlock()
	}
	e.consumerWaitGroup.Done()
}

func Export(cid string) (path string) {
	exporter := NewExporter()
	dirMeta := exporter.ScanDir(cid)
	exportName := fmt.Sprintf("115sha1_%s_%dGB.json", dirMeta.DirName,
		exporter.FileTotalSize>>30)
	outPath, err := dirMeta.Dump(exportName)
	if err != nil {
		fmt.Println("导出到文件失败")
		return
	}

	fmt.Printf("导出文件%s成功, 文件数： %d\n", exportName, exporter.FileCount)
	return outPath
}
