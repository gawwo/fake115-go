package main

import (
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/core"
	"github.com/gawwo/fake115-go/dir"
	"os"
	"testing"
)

// 测试前确定cookie是否已登录
func init() {
	// 确保cookie在登录状态
	loggedIn := core.SetUserInfoConfig()
	if !loggedIn {
		fmt.Println("Login expire or fail...")
		os.Exit(1)
	}
}

// 扫描当前层的文件夹，不涉及下一层文件夹
func TestScanDirWithOffset(t *testing.T) {
	netDir, err := core.ScanDirWithOffset("1932902800904947822", 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(netDir)
}

// 扫描整个文件夹
func TestExport(t *testing.T) {
	config.Debug = true
	config.WorkerNum = 5
	exporter := core.NewExporter()
	dirMeta := exporter.ScanDir("1898007427015248622")
	_, err := dirMeta.Dump("ump_result.txt")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("total size: %dGB, files: %d\n", exporter.FileTotalSize>>30, exporter.FileCount)
}

// 手动查看任务执行情况
func TestImport(t *testing.T) {
	config.Debug = true
	core.Import("1951041685426014426", "ump_result.txt")
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

func TestMakeNetDir(t *testing.T) {
	testDir := dir.Dir{DirName: "for test"}
	testDir.MakeNetDir("0")
}

func TestImportFile(t *testing.T) {
	testFileString := "352Dora.wmv|1447618913|06CCC77F31F4269B5FEB32E7762D0FD" +
		"7C62B1DB9|F29386C9F238CD578BCCAD824FE549F851551473"
	netFile := core.CreateNetFile(testFileString)
	netFile.Cid = "0"
	netFile.Import()
}

func recoverReturn() string {
	// params的顺序不影响defer中对它的读取
	params := false
	defer func() {
		fmt.Println(params)
		if err := recover(); err != nil {
			fmt.Printf("catch error: %v", err)
		}
	}()

	params = true
	panic("normal panic")
}

// 测试错误恢复
func TestRecover(t *testing.T) {
	r := recoverReturn()
	fmt.Println(r)
}
