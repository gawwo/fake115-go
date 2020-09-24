package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSha1(t *testing.T) {
	filePath := "/Users/cheny/Downloads/00416668afe59a1b2f55c0f02a7c277f.mp3"
	expectDigest := "e67bcdfc73ef493197fdd86957e1347fef4eb5de"

	digest, err := FileSha1(filePath)
	if err != nil {
		t.Errorf("摘要生成出错%s", err)
	} else if digest != expectDigest {
		t.Errorf("摘要不匹配, 期望: %s  实际: %s", expectDigest, digest)
	}
}

func TestSha1(t *testing.T) {
	filePath := "/Users/cheny/Downloads/00416668afe59a1b2f55c0f02a7c277f.mp3"
	expectDigest := "e67bcdfc73ef493197fdd86957e1347fef4eb5de"
	f, err := os.Open(filePath)
	if err != nil {
		t.Errorf("找不到文件")
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("读取文件出错")
	}
	digest := Sha1(buf)
	if digest != expectDigest {
		t.Errorf("sha1计算不正确")
	}
}

func BenchmarkFileSha1(b *testing.B) {
	filePath := "/Users/cheny/Downloads/00416668afe59a1b2f55c0f02a7c277f.mp3"

	for i := 0; i < b.N; i++ {
		_, _ = FileSha1(filePath)
	}
}

func BenchmarkSha1(b *testing.B) {
	filePath := "/Users/cheny/Downloads/00416668afe59a1b2f55c0f02a7c277f.mp3"
	f, err := os.Open(filePath)
	if err != nil {
		b.Errorf("找不到文件")
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("读取文件出错")
	}

	for i := 0; i < b.N; i++ {
		Sha1(buf)
	}
}
