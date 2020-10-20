package httpserver

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudger/database"
	uuid "github.com/satori/go.uuid"
)

//QuestionCases 题目数据集
type QuestionCases struct {
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
}

//QuestionInfo 题目信息
type QuestionInfo struct {
	Owner        string        `json:"owner"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	ExampleCases QuestionCases `json:"examplecases"`
}

func handleQuestionInfoRequest(c *gin.Context) {

}

func handleEditQuestionRequest(c *gin.Context) {

}

//AddQuestionRequest 添加题目请求
type AddQuestionRequest struct {
	Token        string        `json:"token" binding:"required"`
	Title        string        `json:"title" binding:"required"`
	Description  string        `json:"description" binding:"required"`
	ExampleCases QuestionCases `json:"examplecases" binding:"required"`
	Cases        QuestionCases `json:"cases" binding:"required"`
}

func handleAddQuestionRequest(c *gin.Context) {
	var aqr AddQuestionRequest
	if c.BindJSON(&aqr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, aqr.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, aqr.Token))
	var qinfo QuestionInfo
	var qcase QuestionCases
	qinfo.Description = aqr.Description
	qinfo.ExampleCases = aqr.ExampleCases
	qinfo.Owner = username
	qinfo.Title = aqr.Title
	qcase = aqr.Cases
	questionID := uuid.NewV4().String()
	qinfoData, _ := json.Marshal(qinfo)
	qcaseData, _ := json.Marshal(qcase)
	database.SetValue(questionDB, questionDescriptionBK, questionID, qinfoData, 0)
	database.SetValue(questionDB, questionCasesBK, questionID, qcaseData, 0)
	database.SAdd(userDB, usernameCreatedQuestionsBK, username, []byte(questionID))
	c.JSON(200, gin.H{"code": 0, "msg": "题目创建成功", "id": questionID})
}

func handleRemoveQuestionRequest(c *gin.Context) {

}

//getInputAndOutputByQuestionName 目前是这样，测试用
func getInputAndOutputByQuestionName(questionName string) (input []string, output []string) {
	input = []string{"hello\n"}
	output = []string{"hello"}
	return
}
