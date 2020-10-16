package golang

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lwh9346/WhaleJudger/iohelper"
)

const imageName = "golang:1.15.2"

func compile(containerName string) {
	containerExec(containerName, "go build -o /root/test.out /root/main.go")
}

func createContainer(containerName string) {
	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	createContainerCMD := exec.Command("docker", "run", "-id", "--name="+containerName, "--net=none", "-v", cwd+"/docker/"+containerName+"/sandbox:/root", imageName)
	err := createContainerCMD.Run()
	if err != nil {
		//Docker error
	}
}

func setUpEnvironmen(containerName, sourceCode string) {
	iohelper.WriteStringToFile("./docker/"+containerName+"/sandbox/main.go", sourceCode)
}

func containerExec(containerName string, command string) error {
	args := append([]string{"exec", "-i", containerName}, strings.Split(command, " ")...)
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

//Debug 调试用
func Debug(containerName string) {
	createContainer(containerName)
	code, _ := ioutil.ReadFile("./main.go")
	setUpEnvironmen(containerName, string(code))
	compile(containerName)
}
