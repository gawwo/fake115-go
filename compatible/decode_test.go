package compatible

import (
	"fmt"
	"github.com/gawwo/fake115-go/dir"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	metaDir := Decode("shoucang.json")
	fmt.Println(metaDir)
}

func TestRebuildTree(t *testing.T) {
	parts := []string{"第一层", "第二层", "第三层"}
	metaDir := &dir.Dir{DirName: "meta", Dirs: []*dir.Dir{{DirName: "第一层"}}}
	last := rebuildTree(metaDir, parts)
	fmt.Println(metaDir)
	fmt.Println(last)
}

func TestFlattenTxtDecode(t *testing.T) {
	file, err := os.Open("ump_result.txt")
	if err != nil {
		println(err.Error())
		println("没有找到文件")
		return
	}
	defer file.Close()

	f := FlattenTxt{}
	metaDir, _ := f.Decode(file)
	fmt.Println(metaDir)
}

func TestJsonFormatDecode(t *testing.T) {
	file, err := os.Open("shoucang.json")
	if err != nil {
		println(err.Error())
		println("没有找到文件")
		return
	}
	defer file.Close()

	j := JsonFormat{}
	metaDir, _ := j.Decode(file)
	fmt.Println(metaDir)
}
