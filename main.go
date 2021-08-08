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
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&config.Debug, "d", false, "Debug mode")
	flag.IntVar(&config.WorkerNum, "n", config.WorkerNum, "worker number")

	// 尝试从外来配置设置cookie文件路径
	flag.StringVar(&config.CookiePath, "cp", config.DefaultCookiePath, "Cookie Path")

	// 尝试从外来配置设置cookie
	flag.StringVar(&config.Cookie, "c", "", "Cookies")
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
