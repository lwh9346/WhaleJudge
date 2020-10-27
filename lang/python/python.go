package python

import (
	"github.com/lwh9346/WhaleJudge/docker"
	"github.com/lwh9346/WhaleJudge/iohelper"
)

//Prepare 进行测试前的各种准备，包括创建容器及编译，返回编译错误信息（如果有）以及测试时的启动参数
func Prepare(sourceCode, containerName string) (errInfo string, runArgs []string) {
	docker.CreateContainer(containerName, "python:3.9.0-alpine3.12")
	iohelper.WriteStringToFile("./docker/"+containerName+"/sandbox/src.py", sourceCode)
	runArgs = []string{"python", "/root/src.py"}
	return "", runArgs
}
