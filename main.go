package main

import (
	"github.com/starjun/gotools/v2"
	"github.com/starjun/toes-cli/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	//config.Home, _ = os.UserHomeDir()
	//fmt.Println("当前执行目录：", config.Home)
	gotools.PathExists("./test")
}
