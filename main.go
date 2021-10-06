package main

import (
	"flag"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/core"
	"github.com/gawwo/fake115-go/log"
	"github.com/gawwo/fake115-go/utils"
	"os"
	"strings"
)

var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "v", false, "显示版本")
	flag.BoolVar(&config.Debug, "d", false, "调试模式")
	flag.IntVar(&config.WorkerNum, "n", config.WorkerNum, "同时进行的数量")
	flag.IntVar(&config.NetworkInterval, "i", config.NetworkInterval, "网络等待间隔")
	flag.IntVar(&config.FilterSize, "f", config.FilterSize, "过滤小于此大小的文件，单位KB")

	// 尝试从外来配置设置cookie文件路径
	flag.StringVar(&config.CookiePath, "cp", config.DefaultCookiePath, "Cookie文件路径")

	// 尝试从外来配置设置cookie
	flag.StringVar(&config.Cookie, "c", "", "Cookie内容")
	// 确保cookie是否真的存在
	if config.Cookie == "" {
		cookie, _ := utils.ReadCookieFile()
		cookie = strings.Trim(cookie, "\n")
		config.Cookie = cookie
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	if showVersion {
		fmt.Println(config.Version)
		return
	}

	if config.Debug {
		config.Logger = log.InitLogger(config.ServerName, true)
	}

	if config.WorkerNum <= 0 {
		config.WorkerNum = 1
	}

	// 确保cookie在登录状态
	loggedIn := core.SetUserInfoConfig()
	if !loggedIn {
		fmt.Println("Login expire or fail...")
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Println("Too few arguments")
		return
	} else if len(args) == 1 {
		core.Export(args[0])
	} else if len(args) == 2 {
		core.Import(args[0], args[1])
	} else {
		fmt.Println("Too much arguments")
		return
	}
}
