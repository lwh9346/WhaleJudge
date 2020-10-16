package golang

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/lwh9346/WhaleJudger/docker"

	"github.com/lwh9346/WhaleJudger/iohelper"
)

func compile(containerName string) string {
	out, _ := containerExec(containerName, "go build -o /root/test.out /root/main.go")
	if out != "" {
		out = fmt.Sprintf("CE:编译错误\n错误信息：\n%s", out)
	}
	return out
}

func setUpEnvironmen(containerName, sourceCode string) {
	iohelper.WriteStringToFile("./docker/"+containerName+"/sandbox/main.go", sourceCode)
}

func containerExec(containerName string, command string) (string, error) {
	args := append([]string{"exec", "-i", containerName}, strings.Split(command, " ")...)
	cmd := exec.Command("docker", args...)
	stderrpipe, err := cmd.StderrPipe()
	reader := bufio.NewReader(stderrpipe)
	cmd.Start()
	outb, _ := ioutil.ReadAll(reader)
	cmd.Wait()
	return string(outb), err
}

//Prepare 进行测试前的各种准备，包括创建容器及编译，返回编译错误信息（如果有）以及测试时的启动参数
func Prepare(sourceCode, containerName string) (errInfo string, runArgs []string) {
	docker.CreateContainer(containerName, "golang:1.15.2")
	setUpEnvironmen(containerName, sourceCode)
	errInfo = compile(containerName)
	runArgs = []string{"/root/test.out"}
	return
}
