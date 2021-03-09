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
	metaDir := &dir.Dir{}
	for _, decoder := range decoders {
		f, err := os.Open(metaPath)
		if err != nil {
			config.Logger.Error("import file not exists",
				zap.String("reason", err.Error()),
				zap.String("path", metaPath))
			fmt.Println("读取导入文件错误")
			return nil
		}

		metaDir, err := decoder.Decode(f)
		if err != nil {
			continue
		}
		if metaDir == nil {
			continue
		}
		if len(metaDir.Files) == 0 || len(metaDir.Dirs) == 0 {
			continue
		}

		f.Close()
		found = true
		break
	}

	if !found {
		return nil
	}
	return metaDir
}
