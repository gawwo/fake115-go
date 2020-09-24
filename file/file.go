package file

import "encoding/json"

type File struct {
	FileName string   `json:"file_name"`
	Files    []string `json:"files"`
	Dirs     []*File  `json:"dirs"`
}

func (file *File) Dump() (string, error) {
	data, err := json.Marshal(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (file *File) Load(fileContent string) (*File, error) {
	err := json.Unmarshal([]byte(fileContent), file)
	if err != nil {
		return nil, err
	}
	return file, nil
}
