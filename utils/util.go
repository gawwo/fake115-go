package utils

import (
	"path/filepath"
)

func FileNameStrip(fullName string) string {
	fileNameBase := filepath.Base(fullName)
	fileSuffix := filepath.Ext(fullName)
	filePrefix := fileNameBase[0 : len(fileNameBase)-len(fileSuffix)]
	return filePrefix
}
