package httpserver

import (
	"fmt"
	"log"

	"github.com/lwh9346/WhaleJudge/lang/python"

	"github.com/lwh9346/WhaleJudge/lang/cpp"

	"github.com/lwh9346/WhaleJudge/database"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudge/docker"
	"github.com/lwh9346/WhaleJudge/judge"
	"github.com/lwh9346/WhaleJudge/lang/golang"
)

//JudgeRequest 一个测评请求
type JudgeRequest struct {
	Language     string `json:"language" binding:"required"`
	SourceCode   string `json:"sourcecode" binding:"required"`
	QuestionName string `json:"question" binding:"required"`
	Token        string `json:"token" binding:"required"`
}

//JudgeResponse 测评结果
type JudgeResponse struct {
	//Msg 包含了每个case的信息，可以是错误信息也可以是其他信息
	Msg []string `json:"msg"`
	//Code 状态码，详见judge.go
	Code int `json:"code"`
}

//TODO:修改返回值风格，加入用户校验
func handleJudgeRequest(c *gin.Context) {
	var request JudgeRequest
	var response JudgeResponse
	if c.BindJSON(&request) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": []string{"请求格式不正确"}})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, request.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": []string{"登陆失败，请重新登陆"}})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, request.Token))
	containerName := GetContainerName()
	var errInfo string
	var args []string //程序运行的参数，这个描述其实不准确，因为程序名也包含在里面
	switch request.Language {
	case "go":
		errInfo, args = golang.Prepare(request.SourceCode, containerName)
		defer docker.KillAndRemoveContainer(containerName)
	case "cpp":
		errInfo, args = cpp.Prepare(request.SourceCode, containerName)
		defer docker.KillAndRemoveContainer(containerName)
	case "python":
		errInfo, args = python.Prepare(request.SourceCode, containerName)
		defer docker.KillAndRemoveContainer(containerName)
	default:
		c.JSON(400, gin.H{"code": 1, "msg": []string{"不支持的语言类型"}})
		return
	}
	if errInfo != "" {
		response.Msg = []string{errInfo}
		response.Code = judge.CompileError
		c.JSON(200, response)
		return
	}
	input, stdOutput := getInputAndOutputByQuestionName(request.QuestionName)
	if len(input) == 0 {
		c.JSON(404, gin.H{"code": 1, "msg": []string{fmt.Sprintf("找不到题目：%s", request.QuestionName)}})
	}
	if len(input) != len(stdOutput) {
		log.Printf("严重错误：题目问题数与答案数不相等\n题目名称：%s\n", request.QuestionName)
		c.JSON(500, gin.H{"code": 1, "msg": []string{"服务器故障，题目错误，请联系管理员"}})
		return
	}
	for k := range input {
		output, code := judge.SingleCase(containerName, input[k], stdOutput[k], args)
		response.Msg = append(response.Msg, output)
		if code != judge.Pass {
			response.Code = code
			c.JSON(200, response)
			return
		}
	}
	//Passed
	database.SAdd(userDB, usernamePassedQuestionsBK, username, []byte(request.QuestionName))
	c.JSON(200, response)
	return
}
