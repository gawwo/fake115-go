package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/gawwo/fake115-go/dir"

	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
	"go.uber.org/zap"
)

// 除非是起始的文件夹，否则其他所有任务都需要先建文件夹，再进行导入工作
func importDir(pid string, meta *dir.Dir, sem *utils.WaitGroupPool) {
	var newest = false
	defer func() {
		if config.Debug {
			fmt.Println("Dir digger on work number: ", sem.Size())
		}

		runtime.Gosched()

		if !newest {
			sem.Done()
		}
	}()

	if sem == nil {
		sem = dir.ProducerWaitGroupPool
		newest = true
	} else {
		// 当达到pool数量上限时，阻塞
		sem.Add()
	}

	var cid string

	// 需要创建一下文件夹
	cid = meta.MakeNetDir(pid)
	if cid == "" {
		config.Logger.Warn("create dir fail",
			zap.String("name", meta.DirName))
		return
	}

	// 提交导入任务到channel中
	for _, fileString := range meta.Files {
		netFile := CreateNetFile(fileString)
		if netFile == nil {
			config.Logger.Warn("error format net file raw content",
				zap.String("content", fileString))
			continue
		}
		netFile.Cid = cid
		task := ImportTask{File: netFile}
		ImportWorkerChannel <- task
	}

	// 处理内层的文件夹
	for _, itemDir := range meta.Dirs {
		if itemDir.HasFile() {
			go importDir(cid, itemDir, sem)
		}
	}
}

func ImportDir(cid string, meta *dir.Dir) {
	// 开启消费者
	config.ConsumerWaitGroup.Add(config.WorkerNum)
	for i := 0; i < config.WorkerNum; i++ {
		go ImportWorker()
	}

	// 开启生产者
	importDir(cid, meta, nil)

	// 等待生产者资源枯竭之后，关闭channel
	dir.ProducerWaitGroupPool.Wait()
	close(ImportWorkerChannel)

	config.ConsumerWaitGroup.Wait()
}

type txtDir struct {
	DirName string   `json:"dir_name"`
	Files   []string `json:"files"`
	//	Dirs    []*txtDir `json:"dirs"`
}

func Import(cid, metaPath string) {
	f, err := os.Open(metaPath)
	if err != nil {
		config.Logger.Error("import file not exists",
			zap.String("reason", err.Error()),
			zap.String("path", metaPath))
		fmt.Println("读取导入文件错误")
		return
	}
	defer f.Close()
	metaBytes, err := ioutil.ReadAll(f)
	if err != nil {
		config.Logger.Error("reader import file error",
			zap.String("reason", err.Error()),
			zap.String("path", metaPath))
		fmt.Println("读取导入文件错误")
		return
	}
	if strings.Index(metaPath, ".txt") != -1 {
		// 开始txt文件目录支持
		var txtToJson txtDir
		txtdirname := metaPath
		reg := regexp.MustCompile(`.*/`)
		txtdirname = reg.ReplaceAllString(txtdirname, "")
		reg = regexp.MustCompile(`.*\\`)
		txtdirname = reg.ReplaceAllString(txtdirname, "")

		txtToJson.DirName = strings.Replace(txtdirname, `.txt`, "", -1)
		file, err := os.Open(metaPath)
		if err != nil {
			println(err.Error())
			println("没有找到" + metaPath)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			txts := scanner.Text()
			if strings.Contains(txts, "115://") { // 如果字符串里面包含了 115:// ，就进行下面的
				fileS := strings.Replace(txts, "115://", "", -1)

				fmt.Println(fileS)
				txtToJson.Files = append([]string{fileS}, txtToJson.Files...)

			} else if strings.Contains(txts, "|") { // 部分sha1 不包含 115：//
				txtToJson.Files = append([]string{txts}, txtToJson.Files...)
			}
		}
		metaBytes, err = json.Marshal(txtToJson)
		if err != nil {
			fmt.Println("JSON ERR:", err)
		}
		fmt.Println(string(metaBytes))

	}
	// 支持 115优化大师导出的json "fold_name":
	stringFold115 := string(metaBytes)
	if strings.Index(stringFold115, "\"fold_name\":") != -1 {
		stringFold115 = strings.Replace(stringFold115, "\"fold_name\":", "\"dir_name\":", -1)
		stringFold115 = strings.Replace(stringFold115, "\"sub_fold\": [", "\"dirs\": [", -1)
		metaBytes = []byte(stringFold115)
	}
	metaDir := dir.NewDir()
	err = metaDir.Load(metaBytes)
	if err != nil {
		config.Logger.Error("import file format error",
			zap.String("reason", err.Error()),
			zap.String("path", metaPath))
		fmt.Println("导入文件格式错误")
		return
	}

	ImportDir(cid, metaDir)

	fmt.Printf("导入文件%dGB，文件数%d\n", config.TotalSize>>30, config.FileCount)
}
