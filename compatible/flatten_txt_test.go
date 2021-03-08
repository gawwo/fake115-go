package compatible

import (
	"fmt"
	"github.com/gawwo/fake115-go/dir"
	"testing"
)

func TestFlattenTxt_Decode(t *testing.T) {
	parts := []string{"第一层", "第二层", "第三层"}
	metaDir := &dir.Dir{DirName: "meta", Dirs: []*dir.Dir{{DirName: "第一层"}}}
	last := rebuildTree(metaDir, parts)
	fmt.Println(metaDir)
	fmt.Println(last)
}
