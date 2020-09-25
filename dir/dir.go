package dir

import "encoding/json"

type Dir struct {
	DirName string   `json:"file_name"`
	Files   []string `json:"files"`
	Dirs    []*Dir   `json:"dirs"`
}

func (dir *Dir) Dump() (string, error) {
	data, err := json.Marshal(dir)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (dir *Dir) Load(fileContent string) (*Dir, error) {
	err := json.Unmarshal([]byte(fileContent), dir)
	if err != nil {
		return nil, err
	}
	return dir, nil
}
