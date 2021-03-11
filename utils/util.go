package utils

import (
	"path/filepath"
)

func ReverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func FileNameStrip(fullName string) string {
	fileNameBase := filepath.Base(fullName)
	fileSuffix := filepath.Ext(fullName)
	filePrefix := fileNameBase[0 : len(fileNameBase)-len(fileSuffix)]
	return filePrefix
}
