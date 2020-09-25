package main

import (
	"flag"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/utils"
)

func init() {
	// 尝试从外来配置设置cookie文件路径
	flag.StringVar(&config.CookiePath, "cp", config.DefaultCookiePath, "Cookie Path")

	// 尝试从外来配置设置cookie
	flag.StringVar(&config.Cookie, "c", "", "Cookies")
	// 确保cookie是否真的存在
	if config.Cookie == "" {
		cookie, _ := utils.ReadCookieFile()
		config.Cookie = cookie
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	// 没cookie什么都干不了
	if config.Cookie == "" {
		fmt.Println("Load cookie error! Please add cookie or add cookie file path")
		return
	}

	if len(args) < 1 {
		fmt.Println("Too few arguments")
		return
	} else if len(args) == 1 {
		// 导出模式
	} else if len(args) == 2 {
		// 导入模式
	} else {
		fmt.Println("Too much arguments")
		return
	}
}
