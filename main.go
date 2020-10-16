package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/lwh9346/WhaleJudger/httpserver"

	"github.com/lwh9346/WhaleJudger/judge"

	"github.com/lwh9346/WhaleJudger/docker"

	"github.com/lwh9346/WhaleJudger/lang/golang"
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
		debug("codingTest")
	default:
		//default
	}
}

func debug(name string) {
	docker.CreateContainer(name, "golang:1.15.2")
	sourceCode, _ := ioutil.ReadFile("./main.go")
	req := httpserver.JudgeRequest{SourceCode: string(sourceCode), Language: "go", QuestionName: "silly question"}
	b, _ := json.Marshal(req)
	ioutil.WriteFile("./test.json", b, 0666)
	errInfo, runArgs := golang.Prepare(string(sourceCode), name)
	if errInfo != "" {
		log.Println(errInfo)
		return
	}
	output, _ := judge.SingleCase(name, "23 14\n", "37", runArgs)
	log.Println(output)
}
