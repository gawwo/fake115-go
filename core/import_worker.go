package core

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"go.uber.org/zap"
	"time"
)

type ImportTask struct {
	File *NetFile
}

var ImportWorkerChannel = make(chan ImportTask, config.WorkerNum*config.WorkerNumRate)

func ImportWorker() {
	for task := range ImportWorkerChannel {
		if config.Debug {
			fmt.Println("channel len: ", len(ImportWorkerChannel))
		}

		start := time.Now().Unix()
		result := task.File.Import()
		if !result {
			config.Logger.Warn("import failed", zap.String("name", task.File.Name))
			continue
		}

		elapsed := time.Now().Unix() - start
		if elapsed > int64(3) {
			config.Logger.Warn("task slow", zap.String("name", task.File.Name),
				zap.Int64("elapsed", elapsed))
		}

		lock.Lock()
		config.TotalSize += task.File.Size
		config.FileCount += 1
		lock.Unlock()
	}
	config.ConsumerWaitGroup.Done()
}
