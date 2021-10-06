package core

import (
	"fmt"
	"github.com/gawwo/fake115-go/compatible"
	"runtime"
	"sync"
	"time"

	"github.com/gawwo/fake115-go/dir"

	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
)

type ImportTask struct {
	File *NetFile
}

type importer struct {
	taskChannel       chan ImportTask
	consumerWaitGroup sync.WaitGroup
	// 通过pool支持设置上限
	producerWaitGroupPool *utils.WaitGroupPool
	lock                  sync.Mutex
	FileCount             int
	FileTotalSize         int
}

func NewImporter() *importer {
	return &importer{
		taskChannel:           make(chan ImportTask, config.WorkerNum*config.WorkerNumRate),
		consumerWaitGroup:     sync.WaitGroup{},
		producerWaitGroupPool: utils.NewWaitGroupPool(config.WorkerNum),
		FileCount:             0,
		FileTotalSize:         0,
	}
}

func (i *importer) importDir(pid string, meta *dir.Dir) {
	defer func() {
		if config.Debug {
			fmt.Println("Dir digger on work number: ", i.producerWaitGroupPool.Size())
		}

		runtime.Gosched()
		i.producerWaitGroupPool.Done()
	}()

	i.producerWaitGroupPool.Add()

	time.Sleep(time.Second * time.Duration(config.NetworkInterval))

	var cid string

	// 需要创建一下文件夹
	cid = meta.MakeNetDir(pid)
	if cid == "" {
		config.Logger.Warn("create dir fail",
			zap.String("name", meta.DirName))
		return
	}

	// 提交导入任务到channel中
	for _, fileString := range meta.Files {
		netFile := CreateNetFile(fileString)
		if netFile == nil {
			config.Logger.Warn("error format net file raw content",
				zap.String("content", fileString))
			continue
		}
		// 过滤太小的文件
		if config.FilterSize<<10 > netFile.Size {
			continue
		}
		netFile.Cid = cid
		task := ImportTask{File: netFile}
		i.taskChannel <- task
	}

	// 处理内层的文件夹
	for _, itemDir := range meta.Dirs {
		if itemDir.HasFile() {
			go i.importDir(cid, itemDir)
		}
	}
}

func (i *importer) importConsumer() {
	for task := range i.taskChannel {
		if config.Debug {
			fmt.Println("channel len: ", len(i.taskChannel))
		}

		time.Sleep(time.Second * time.Duration(config.NetworkInterval))
		result := task.File.Import()
		if !result {
			config.Logger.Warn("import failed", zap.String("name", task.File.Name))
			continue
		}

		i.lock.Lock()
		i.FileTotalSize += task.File.Size
		i.FileCount += 1
		i.lock.Unlock()
	}
	i.consumerWaitGroup.Done()
}

func (i *importer) ImportDir(cid string, meta *dir.Dir) {
	// 开启消费者
	i.consumerWaitGroup.Add(config.WorkerNum)
	for n := 0; n < config.WorkerNum; n++ {
		go i.importConsumer()
	}

	// 开启生产者
	i.importDir(cid, meta)

	// 等待生产者资源枯竭之后，关闭channel
	i.producerWaitGroupPool.Wait()
	close(i.taskChannel)

	i.consumerWaitGroup.Wait()
}

func Import(cid, metaPath string) {
	metaDir := compatible.Decode(metaPath)
	if metaDir == nil {
		fmt.Println("未找到可处理的格式")
		return
	}

	importer := NewImporter()
	importer.ImportDir(cid, metaDir)

	fmt.Printf("导入文件%dGB，文件数%d\n", importer.FileTotalSize>>30, importer.FileCount)
}
