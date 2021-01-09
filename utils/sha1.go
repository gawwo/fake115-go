package utils

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func FileSha1(filePath string) (hexDigest string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return hexDigest, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	buf := make([]byte, 64<<10)
	sha1Hash := sha1.New()

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				goto readFinish
			}
			return hexDigest, err
		}

		sha1Hash.Write(buf[:n])
	}

readFinish:
	hexDigest = hex.EncodeToString(sha1Hash.Sum(nil))
	return hexDigest, nil
}

func Sha1(content []byte) string {
	sha1Hash := sha1.New()
	sha1Hash.Write(content)
	return hex.EncodeToString(sha1Hash.Sum(nil))
}
