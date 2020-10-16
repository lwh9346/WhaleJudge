package golang

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lwh9346/WhaleJudger/iohelper"
)

const imageName = "golang:1.15.2"

func compile(containerName string) {
	containerExec(containerName, "cd /root && go build -o test.out")
}

func createContainer(containerName string) {
	cwd := filepath.Dir(os.Args[0])
	createContainerCMD := exec.Command("docker", "run", "-id", "--name=\""+containerName+"\"", "--net=none", "-v", cwd+"/docker/"+containerName+"/sandbox:/root", imageName)
	err := createContainerCMD.Run()
	if err != nil {
		//Docker error
	}
}

func setUpEnvironmen(containerName string) {
	containerExec(containerName, "cd /root && go mod init test")
	iohelper.CopyFile("./main.go", "./docker/"+containerName+"/sandbox/main.go")
}

func containerExec(containerName string, command string) error {
	cmd := exec.Command("docker", "exec", "-i", containerName, command)
	return cmd.Run()
}

//Debug 调试用
func Debug(containerName string) {
	createContainer(containerName)
	setUpEnvironmen(containerName)
	compile(containerName)
}
