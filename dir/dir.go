package dir

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
	"os"
)

type Dir struct {
	DirName string   `json:"dir_name"`
	Files   []string `json:"files"`
	Dirs    []*Dir   `json:"dirs"`
}

type dirMake struct {
	State bool   `json:"state"`
	Error string `json:"error"`
	Cid   string `json:"cid"`
}

// 强行指定初始化，防止json之后，Dirs和Files为null
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

func (dir *Dir) Load(fileContent []byte) error {
	err := json.Unmarshal(fileContent, dir)
	if err != nil {
		return err
	}
	return nil
}

// 递归探测文件夹中是否有文件
func (dir *Dir) HasFile() bool {
	if len(dir.Files) > 0 {
		return true
	}

	for _, innerDir := range dir.Dirs {
		if innerDir.HasFile() {
			return true
		}
	}
	return false
}

// 创建新的文件夹
func (dir *Dir) MakeNetDir(pid string) string {
	defer func() {
		if err := recover(); err != nil {
			config.Logger.Error("create dir error",
				zap.String("name", dir.DirName),
				zap.String("reason", fmt.Sprintf("%v", err)))
		}
	}()

	url := "https://webapi.115.com/files/add"
	for i := 0; i < config.RetryTimes; i++ {
		var dirName string
		if i != 0 {
			dirName = fmt.Sprintf("%s_%d", dir.DirName, i)
		} else {
			dirName = dir.DirName
		}

		data := map[string]string{
			"pid":   pid,
			"cname": dirName,
		}
		headers := config.GetFakeHeaders(true)
		body, err := utils.PostForm(url, headers, data)
		if err != nil {
			config.Logger.Warn("make dir fail", zap.String("name", dirName))
			return ""
		}

		dirMakeResult := new(dirMake)
		err = json.Unmarshal(body, dirMakeResult)
		if err != nil {
			config.Logger.Warn("parse make dir result fail",
				zap.String("reason", err.Error()),
				zap.String("name", dir.DirName))
			return ""
		}

		if dirMakeResult.State {
			return dirMakeResult.Cid
		}

		if dirMakeResult.Error == "该目录名称已存在。" {
			config.Logger.Warn("dir had exists, change name to continue",
				zap.String("name", dir.DirName))
			continue
		}

		return ""
	}
	return ""
}
