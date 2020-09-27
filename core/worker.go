package core

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
)

type Task struct {
	Dir  *dir.Dir
	File *NetFile
}

func Worker() {
	// 知道WorkerChannel关闭前，都一直工作
	for task := range config.WorkerChannel {
		task.File.Export(task.Dir)
	}
	config.WaitGroup.Done()
}
