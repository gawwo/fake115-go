package main

import (
	"fmt"
	"github.com/gawwo/fake115-go/core"
	"github.com/gawwo/fake115-go/dir"
	"testing"
)

func TestScanDirWithOffset(t *testing.T) {
	netDir, err := core.ScanDirWithOffset("1932902800904947822", 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(netDir)
}

func TestNetFileExport(t *testing.T) {
	netFile := core.NetFile{Pc: "awzzj1xe4id8sht7o1"}
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
