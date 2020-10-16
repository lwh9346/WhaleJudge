package golang

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lwh9346/WhaleJudger/judge"

	"github.com/lwh9346/WhaleJudger/iohelper"
)

const imageName = "golang:1.15.2"

func compile(containerName string) {
	_, err := containerExec(containerName, "go build -o /root/test.out /root/main.go")
	if err != nil {
		log.Printf("编译错误\n")
	}
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

func containerExec(containerName string, command string) ([]byte, error) {
	args := append([]string{"exec", "-i", containerName}, strings.Split(command, " ")...)
	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	return out, err
}

//Debug 调试用
func Debug(containerName string) {
	createContainer(containerName)
	code, _ := ioutil.ReadFile("./main.go")
	setUpEnvironmen(containerName, string(code))
	compile(containerName)
	//out, _ := containerExec(containerName, "/root/test.out")
	outstr, _ := judge.SingleCase(containerName, "hello\n", "hello", []string{"/root/test.out"})
	log.Println(outstr)
}
