package dir

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	info := `{"file_name": "\u4e16\u754c\u81ea\u7136\u9057\u4ea7", "files": ["readme.txt"], "dirs": [{"file_name": "\u4e2d\u56fd", "files": ["\u6210\u90fd.mp4"], "dirs": []}]}`
	file := NewDir()
	err := file.Load([]byte(info))
	if err != nil {
		t.Error("Load error", err.Error())
	}
}

func TestDump(t *testing.T) {
	mark := "中国"
	file := Dir{
		DirName: "mark",
		Files:   []string{"readme.txt"},
		Dirs: []*Dir{
			{
				DirName: "中国",
				Files:   []string{"成都.mp4"},
				Dirs:    []*Dir{},
			},
		},
	}

	dump, err := file.Dumps()
	if err != nil {
		t.Errorf(err.Error())
	}
	contains := strings.Contains(string(dump), mark)
	if !contains {
		t.Errorf("format not contain mark: %s", dump)
	}
}
