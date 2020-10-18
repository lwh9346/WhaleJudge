package main

import (
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
		dataBaseBasicFunctionTest()
	default:
		//default
	}
}
