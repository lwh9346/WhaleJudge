package httpserver

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudger/docker"
	"github.com/lwh9346/WhaleJudger/judge"
	"github.com/lwh9346/WhaleJudger/lang/golang"
)

func handleJudgeRequest(c *gin.Context) {
	var request JudgeRequest
	var response JudgeResponse
	if c.BindJSON(&request) != nil {
		c.JSON(400, gin.H{"err": "请求格式不正确"})
		return
	}
	containerName := getContainerName() //这只是测试用的，正式服务器里面需要获取一个唯一的名称
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
	if len(input) != len(stdOutput) {
		log.Printf("严重错误：题目问题数与答案数不相等\n题目名称：%s\n", request.QuestionName)
		c.JSON(400, gin.H{"err": "服务器故障：题目错误"})
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
