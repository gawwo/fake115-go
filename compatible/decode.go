package compatible

import (
	"github.com/gawwo/fake115-go/dir"
	"os"
)

const flattenTxtPrefix = "115://"
const flattenTxtSplit = "|"
const normalSplitLen = 4

type Decoder interface {
	Decode(file *os.File) (*dir.Dir, error)
}
