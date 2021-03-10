package compatible

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"go.uber.org/zap"
	"os"
)

const flattenTxtPrefix = "115://"
const flattenTxtSplit = "|"
const normalSplitLen = 4

type Decoder interface {
	Decode(file *os.File) (*dir.Dir, error)
}

func Decode(metaPath string) *dir.Dir {
	decoders := []Decoder{
		&SelfJson{},
		&FlattenTxt{},
		&JsonFormat{},
	}
	found := false
	var metaDir = &dir.Dir{}
	for _, decoder := range decoders {
		f, err := os.Open(metaPath)
		if err != nil {
			config.Logger.Error("import file not exists",
				zap.String("reason", err.Error()),
				zap.String("path", metaPath))
			fmt.Println("读取导入文件错误")
			return nil
		}

		decodeDir, err := decoder.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		if decodeDir == nil {
			continue
		}
		if len(decodeDir.Files) == 0 && len(decodeDir.Dirs) == 0 {
			continue
		}

		found = true
		metaDir = decodeDir
		break
	}

	if !found {
		return nil
	}
	return metaDir
}
