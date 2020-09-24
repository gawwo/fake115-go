package file

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	info := `{"file_name": "\u4e16\u754c\u81ea\u7136\u9057\u4ea7", "files": ["readme.txt"], "dirs": [{"file_name": "\u4e2d\u56fd", "files": ["\u6210\u90fd.mp4"], "dirs": []}]}`
	//file := new(File)
	var file File
	fileObj, _ := file.Load(info)
	if fileObj == nil {
		t.Errorf("Load error")
	}
}

func TestDump(t *testing.T) {
	mark := "世界自然遗产"
	file := File{
		FileName: "mark",
		Files:    []string{"readme.txt"},
		Dirs: []*File{
			{
				FileName: "中国",
				Files:    []string{"成都.mp4"},
				Dirs:     []*File{},
			},
		},
	}

	dump, err := file.Dump()
	if err != nil {
		t.Errorf(err.Error())
	}
	contains := strings.Contains(dump, mark)
	if !contains {
		t.Errorf("format not contain mark: %s", dump)
	}
}
