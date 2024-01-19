package handler

import (
	"log"
	"os"
)

const defaultPort = 8000

var NginxConfig string

func init() {
	if len(os.Args) <= 1 {
		log.Println("[info] 没有指定nginx配置文件路径")
		os.Exit(-1)
	}
	NginxConfig = os.Args[1]
	_, err := parse(NginxConfig)
	if err != nil {
		log.Println("[info] nginx配置文件解析失败")
		os.Exit(-1)
	}
}

func ParseArgPort() int {
	return defaultPort
}
