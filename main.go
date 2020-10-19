package main

import (
	"log"
	"os"

	"github.com/lwh9346/WhaleJudger/httpserver"
)

func main() {
	if len(os.Args) == 1 {
		httpserver.StartHTTPServer()
		return
	}
	switch os.Args[1] {
	case "init":
		//init
	case "test":
		if len(os.Args) < 3 {
			log.Fatalln("test命令没有足够的参数")
		}
		switch os.Args[2] {
		case "docker":
			dockerBasicFunctionTest()
		case "database":
			dataBaseBasicFunctionTest()
		}
	default:
		//default
	}
}
