package compatible

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type SelfJson struct{}

func (s *SelfJson) Decode(file *os.File) (*dir.Dir, error) {
	metaBytes, err := ioutil.ReadAll(file)
	if err != nil {
		config.Logger.Error("import file not exists",
			zap.String("reason", err.Error()),
			zap.String("path", file.Name()))
		return nil, err
	}

	metaDir := dir.NewDir()
	err = metaDir.Load(metaBytes)
	if err != nil {
		return nil, err
	} else {
		return metaDir, nil
	}
}
