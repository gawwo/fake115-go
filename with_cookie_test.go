package main

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/core"
	"testing"
)

// 扫描当前层的文件夹，不涉及下一层文件夹
func TestScanDirWithOffset(t *testing.T) {
	netDir, err := core.ScanDirWithOffset("1932902800904947822", 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(netDir)
}

// 扫描整个文件夹
func TestScanDir(t *testing.T) {
	config.Debug = true
	dir := core.ScanDir("1898007427015248622")
	_, err := dir.Dump("ump_result.json")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("total size: %dGB\n", config.TotalSize>>30)
}

func TestNetFileExport(t *testing.T) {
	netFile := core.NetFile{
		Fid:  "1932902801198549107",
		Cid:  "1932902800904947822",
		Size: 3153756278,
		Name: "raised.by.wolves.2020.s01e07.1080p.web.h264-videohole.mkv",
		Sha:  "44451C2DDCE125722FBA9DE1760E55E265023A73",
		Pc:   "b9zzwuk9729f283dt",
	}
	result := netFile.Export()
	if result == "" {
		t.Error("export file fail")
	}
}

func TestRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch error: %v", err)
		}
	}()

	panic("panic once")
}
