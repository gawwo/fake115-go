package dir

import (
	"bufio"
	"encoding/json"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"os"
)

// 扫描文件夹时，用于控制扫描goroutine数量的信号量池
var WaitGroupPool = utils.NewWaitGroupPool(config.WorkerNum)

type Dir struct {
	DirName string   `json:"file_name"`
	Files   []string `json:"files"`
	Dirs    []*Dir   `json:"dirs"`
}

func NewDir() *Dir {
	return &Dir{Dirs: []*Dir{}, Files: []string{}}
}

func (dir *Dir) Dumps() ([]byte, error) {
	data, err := json.Marshal(dir)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 导出到文件中
func (dir *Dir) Dump(outPath string) (string, error) {
	data, err := dir.Dumps()
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	_, err = writer.Write(data)
	if err != nil {
		return "", nil
	}
	writer.Flush()

	return outPath, nil
}

func (dir *Dir) Load(fileContent string) (*Dir, error) {
	err := json.Unmarshal([]byte(fileContent), dir)
	if err != nil {
		return nil, err
	}
	return dir, nil
}
