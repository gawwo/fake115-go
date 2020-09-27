package main

import (
	"fmt"
	"github.com/gawwo/fake115-go/dir"
	"github.com/gawwo/fake115-go/executor"
	"testing"
)

func TestScanDirWithOffset(t *testing.T) {
	executor.ScanDirWithOffset("1932902800904947822", 0)
}

func TestNetFileExport(t *testing.T) {
	netFile := executor.NetFile{Pc: "awzzj1xe4id8sht7o1"}
	dirObj := new(dir.Dir)
	netFile.Export(dirObj)
}

func TestRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch error: %v", err)
		}
	}()

	panic("panic once")
}
