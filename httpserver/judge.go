package httpserver

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudger/docker"
	"github.com/lwh9346/WhaleJudger/judge"
	"github.com/lwh9346/WhaleJudger/lang/golang"
)

//JudgeRequest 一个测评请求
type JudgeRequest struct {
	//Language 程序所使用的语言
	Language string `json:"language" binding:"required"`
	//SourceCode 程序源代码
	SourceCode string `json:"sourcecode" binding:"required"`
	//QuestionName 所回答的问题的名称（ID）
	QuestionName string `json:"question" binding:"required"`
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
		c.JSON(400, gin.H{"err": "请求格式不正确"})
		return
	}
	containerName := GetContainerName()
	var errInfo string
	var args []string //程序运行的参数，这个描述其实不准确，因为程序名也包含在里面
	switch request.Language {
	case "go":
		errInfo, args = golang.Prepare(request.SourceCode, containerName)
		defer docker.KillAndRemoveContainer(containerName)
	default:
		c.JSON(400, gin.H{"err": "不支持的语言类型"})
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
		c.JSON(404, gin.H{"err": fmt.Sprintf("找不到题目：%s", request.QuestionName)})
	}
	if len(input) != len(stdOutput) {
		log.Printf("严重错误：题目问题数与答案数不相等\n题目名称：%s\n", request.QuestionName)
		c.JSON(400, gin.H{"err": "服务器故障，题目错误，请联系管理员"})
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
	c.JSON(200, response)
	return
}
