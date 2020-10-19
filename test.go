package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/xujiajun/nutsdb"

	"github.com/lwh9346/WhaleJudger/docker"
	"github.com/lwh9346/WhaleJudger/httpserver"
	"github.com/lwh9346/WhaleJudger/judge"
	"github.com/lwh9346/WhaleJudger/lang/golang"
)

//这个包的作用是进行各个模块的测试

//测试docker的容器创建、代码编译运行以及容器销毁功能
func dockerBasicFunctionTest() {
	name := httpserver.GetContainerName()
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
	docker.KillAndRemoveContainer(name)
	log.Println(output)
}

func dataBaseBasicFunctionTest() {
	option := nutsdb.DefaultOptions
	option.Dir = "./db/test"
	db, err := nutsdb.Open(option)
	if err != nil {
		log.Fatal(err)
	}

	k := []byte("testkey")
	v := []byte("testvalue")
	var vRead []byte
	bucket := "defaultBucket"

	if err = db.Update(
		func(tx *nutsdb.Tx) error {
			return tx.Put(bucket, k, v, nutsdb.Persistent)
		}); err != nil {
		log.Fatal(err)
	}
	if err = db.View(
		func(tx *nutsdb.Tx) error {
			e, err1 := tx.Get(bucket, k)
			vRead = e.Value
			return err1
		}); err != nil {
		log.Fatal(err)
	}
	log.Println(string(vRead))
	db.Close()
	err = os.RemoveAll("./db/test")
}
