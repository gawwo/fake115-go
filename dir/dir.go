package dir

import (
	"bufio"
	"encoding/json"
	"os"
)

type Dir struct {
	DirName string   `json:"file_name"`
	Files   []string `json:"files"`
	Dirs    []*Dir   `json:"dirs"`
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

	f, err := os.Open(outPath)
	if err != nil {
		return "", err
	}

	writer := bufio.NewWriter(f)
	_, err = writer.Write(data)
	if err != nil {
		return "", nil
	}
	return outPath, nil
}

func (dir *Dir) Load(fileContent string) (*Dir, error) {
	err := json.Unmarshal([]byte(fileContent), dir)
	if err != nil {
		return nil, err
	}
	return dir, nil
}
