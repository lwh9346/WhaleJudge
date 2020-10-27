package httpserver

import (
	"os/exec"
	"time"

	"github.com/lwh9346/WhaleJudge/lang/cpp"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudge/docker"
	"github.com/lwh9346/WhaleJudge/lang/golang"
)

type demoRequest struct {
	SourceCode string `json:"src" binding:"required"`
	Language   string `json:"lang" binding:"required"`
}

func handleDemoRequest(c *gin.Context) {
	var dr demoRequest
	if c.BindJSON(&dr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	var errInfo string
	var args []string
	containerName := GetContainerName()
	switch dr.Language {
	case "go":
		errInfo, args = golang.Prepare(dr.SourceCode, containerName)
		defer docker.KillAndRemoveContainer(containerName)
	case "cpp":
		errInfo, args = cpp.Prepare(dr.SourceCode, containerName)
		defer docker.KillAndRemoveContainer(containerName)
	default:
		c.JSON(200, gin.H{"code": 1, "msg": "不支持的语言类型"})
		return
	}
	if errInfo != "" {
		c.JSON(200, gin.H{"code": 2, "msg": errInfo})
		return
	}
	args = append([]string{"exec", "-i", containerName}, args...)
	cmd := exec.Command("docker", args...)
	outputChan := make(chan string)
	errChan := make(chan error)
	timeOutChan := time.NewTimer(time.Second).C
	go func() {
		o, e := cmd.Output()
		if e != nil {
			errChan <- e
			return
		}
		outputChan <- string(o)
	}()
	select {
	case output := <-outputChan:
		c.JSON(200, gin.H{"code": 0, "msg": output})
		return
	case err := <-errChan:
		c.JSON(200, gin.H{"code": 3, "msg": err.Error()})
		return
	case <-timeOutChan:
		cmd.Process.Kill()
		c.JSON(200, gin.H{"code": 4, "msg": "Time out!"})
		return
	}
}
