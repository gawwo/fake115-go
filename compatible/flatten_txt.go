package compatible

import (
	"bufio"
	"errors"
	"github.com/gawwo/fake115-go/dir"
	"os"
	"strings"
)

type FlattenTxt struct{}

func (f *FlattenTxt) Decode(file *os.File) (*dir.Dir, error) {
	metaDir := &dir.Dir{DirName: file.Name()}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 部分sha1 不包含 115：//
		if strings.HasPrefix(line, flattenTxtPrefix) {
			line = line[len(flattenTxtPrefix):]
		}

		parts := strings.Split(line, flattenTxtSplit)
		if len(parts) < normalSplitLen {
			return metaDir, errors.New("Format Error ")
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

func rebuildTree(metaDir *dir.Dir, paths []string) *dir.Dir {
	// 根据路径重建目录，并返回最后一层目录
	if len(paths) == 0 {
		return metaDir
	}

	found := false
	expectDir := &dir.Dir{}
	for _, innerDir := range metaDir.Dirs {
		if innerDir.DirName == paths[0] {
			found = true
			expectDir = innerDir
		}
	}

	if !found {
		expectDir.DirName = paths[0]
		metaDir.Dirs = append(metaDir.Dirs, expectDir)
	}

	return rebuildTree(expectDir, paths[1:])
}
