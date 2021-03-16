package compatible

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/utils"
)

type FlattenTxt struct{}

func (f *FlattenTxt) Decode(file *os.File) (*dir.Dir, error) {
	metaDir := &dir.Dir{DirName: utils.FileNameStrip(file.Name())}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, flattenTxtPrefix) {
			line = line[len(flattenTxtPrefix):]
		}

		parts := strings.Split(line, flattenTxtSplit)
		if len(parts) < normalSplitLen {
			fmt.Printf("本行字符串有误 %s ", line)
			continue

		} else if len(parts) == normalSplitLen {
			metaDir.Files = append(metaDir.Files, line)

		} else {
			dirParts := parts[normalSplitLen:]
			treeNode := rebuildTree(metaDir, dirParts)
			treeNode.Files = append(treeNode.Files, strings.Join(parts[:normalSplitLen], flattenTxtSplit))
		}
	}
	return metaDir, nil
}

func rebuildTree(metaDir *dir.Dir, dirpaths []string) *dir.Dir {
	if len(dirpaths) == 0 {
		return metaDir
	}

	found := false
	expectDir := &dir.Dir{}
	for _, innerDir := range metaDir.Dirs {
		if innerDir.DirName == dirpaths[0] {
			found = true
			expectDir = innerDir
		}
	}

	if !found {
		expectDir.DirName = dirpaths[0]
		metaDir.Dirs = append(metaDir.Dirs, expectDir)
	}

	return rebuildTree(expectDir, dirpaths[1:])
}
