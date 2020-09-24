package benchmark

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func BufIo(f *os.File) ([]byte, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()

	//r := bufio.NewReaderSize(f, 512 << 10)
	r := bufio.NewReader(f)
	buf := make([]byte, 64<<10)

	content := make([]byte, fileSize)

	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		content = append(content, buf[:n]...)
	}
	return content, nil
}

//
func BenchmarkBufIo(b *testing.B) {
	filePath := "/Users/cheny/OneDrive/Documents/fake115/log/logs/info.log"

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = BufIo(f)
	}
}
