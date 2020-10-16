package docker

import (
	"os"
	"os/exec"
	"path/filepath"
)

//CreateContainer 用指定名称创建Docker容器
func CreateContainer(containerName, imageName string) error {
	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	createContainerCMD := exec.Command("docker", "run", "-id", "--name="+containerName, "--net=none", "-v", cwd+"/docker/"+containerName+"/sandbox:/root", imageName)
	err := createContainerCMD.Run()
	return err
}
