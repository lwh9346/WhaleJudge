package httpserver

import (
	"github.com/gin-gonic/gin"
)

func handleQuestionInfoRequest(c *gin.Context) {

}

func handleEditQuestionRequest(c *gin.Context) {

}

func handleRemoveQuestionRequest(c *gin.Context) {

}

//getInputAndOutputByQuestionName 目前是这样，测试用
func getInputAndOutputByQuestionName(questionName string) (input []string, output []string) {
	input = []string{"hello\n"}
	output = []string{"hello"}
	return
}
