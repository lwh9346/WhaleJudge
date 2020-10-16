package main

import (
	"os"

	"github.com/lwh9346/WhaleJudger/lang/golang"
)

func main() {
	if len(os.Args) == 1 {
		//RunHttpServer
		return
	}
	switch os.Args[1] {
	case "init":
		//init
	case "test":
		golang.Debug("codingTest")
	default:
		//default
	}
}
