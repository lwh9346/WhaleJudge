package cpp

import (
	"bufio"
	"io/ioutil"
	"os/exec"

	"github.com/lwh9346/WhaleJudge/docker"
	"github.com/lwh9346/WhaleJudge/iohelper"
)

//Prepare 进行测试前的各种准备，包括创建容器及编译，返回编译错误信息（如果有）以及测试时的启动参数
func Prepare(sourceCode, containerName string) (errInfo string, runArgs []string) {
	docker.CreateContainer(containerName, "ubuntu:20.04")
	iohelper.WriteStringToFile("./docker/"+containerName+"/src.cpp", sourceCode)
	compileCmd := exec.Command("g++", "./docker/"+containerName+"/src.cpp", "-o", "./docker/"+containerName+"/sandbox/test.out")
	stderrpipe, _ := compileCmd.StderrPipe()
	reader := bufio.NewReader(stderrpipe)
	compileCmd.Start()
	outb, _ := ioutil.ReadAll(reader)
	compileCmd.Wait()
	runArgs = []string{"/root/test.out"}
	return string(outb), runArgs
}
