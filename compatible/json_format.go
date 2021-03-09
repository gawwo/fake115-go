package compatible

import (
	"errors"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

type JsonFormat struct{}

func (j *JsonFormat) Decode(file *os.File) (*dir.Dir, error) {
	metaBytes, err := ioutil.ReadAll(file)
	if err != nil {
		config.Logger.Error("import file not exists",
			zap.String("reason", err.Error()),
			zap.String("path", file.Name()))
		return nil, err
	}

	// 支持 115优化大师导出的json "fold_name":
	stringFold115 := string(metaBytes)
	if strings.Index(stringFold115, "\"fold_name\":") != -1 {
		stringFold115 = strings.Replace(stringFold115, "\"fold_name\":", "\"dir_name\":", -1)
		stringFold115 = strings.Replace(stringFold115, "\"sub_fold\": [", "\"dirs\": [", -1)
		metaBytes = []byte(stringFold115)
	} else {
		return nil, errors.New("Error Format ")
	}

	metaDir := dir.NewDir()
	err = metaDir.Load(metaBytes)
	if err != nil {
		return nil, err
	} else {
		return metaDir, nil
	}
}
