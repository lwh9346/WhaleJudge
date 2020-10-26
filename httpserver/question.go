package httpserver

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudge/database"
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

//QuestionInfoRequest 题目信息请求
type QuestionInfoRequest struct {
	QuestionName string `json:"questionname"`
}

func handleQuestionInfoRequest(c *gin.Context) {
	var qir QuestionInfoRequest
	if c.BindJSON(&qir) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(questionDB, questionDescriptionBK, qir.QuestionName) {
		c.JSON(404, gin.H{"code": 1, "msg": "题目不存在"})
		return
	}
	var qi QuestionInfo
	data := database.GetValue(questionDB, questionDescriptionBK, qir.QuestionName)
	json.Unmarshal(data, &qi)
	c.JSON(200, gin.H{"code": 0, "msg": "查询成功", "info": qi})
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
	if len(aqr.Title) > 64 {
		c.JSON(400, gin.H{"code": 1, "msg": "标题过长，最长64字符"})
		return
	}
	questionName := aqr.Title
	if database.HasKey(questionDB, questionDescriptionBK, questionName) {
		c.JSON(400, gin.H{"code": 1, "msg": "该题目名已存在"})
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
	qinfoData, _ := json.Marshal(qinfo)
	qcaseData, _ := json.Marshal(qcase)
	database.SetValue(questionDB, questionDescriptionBK, questionName, qinfoData, 0)
	database.SetValue(questionDB, questionCasesBK, questionName, qcaseData, 0)
	database.SAdd(userDB, usernameCreatedQuestionsBK, username, []byte(questionName))
	c.JSON(200, gin.H{"code": 0, "msg": "题目创建成功", "name": questionName})
}

//RemoveQuestionRequest 移除某问题的请求
type RemoveQuestionRequest struct {
	Token        string `json:"token" binding:"required"`
	QuestionName string `json:"questionname" binding:"required"`
}

func handleRemoveQuestionRequest(c *gin.Context) {
	var rqr RemoveQuestionRequest
	if c.BindJSON(&rqr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, rqr.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, rqr.Token))
	if !database.SHasKey(userDB, usernameCreatedQuestionsBK, username) {
		c.JSON(400, gin.H{"code": 1, "msg": "你没有删除该问题的权限"})
		return
	}
	if !database.SIsMember(userDB, usernameCreatedQuestionsBK, username, []byte(rqr.QuestionName)) {
		c.JSON(400, gin.H{"code": 1, "msg": "你没有删除该问题的权限"})
	}
	database.SRemove(userDB, usernameCreatedQuestionsBK, username, []byte(rqr.QuestionName))
	database.RemoveKey(questionDB, questionCasesBK, rqr.QuestionName)
	database.RemoveKey(questionDB, questionDescriptionBK, rqr.QuestionName)
	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}

//getInputAndOutputByQuestionName 获取某个问题的输入输出
func getInputAndOutputByQuestionName(questionName string) (input []string, output []string) {
	if !database.HasKey(questionDB, questionCasesBK, questionName) {
		return
	}
	data := database.GetValue(questionDB, questionCasesBK, questionName)
	var questionCases QuestionCases
	json.Unmarshal(data, &questionCases)
	input = questionCases.Inputs
	output = questionCases.Outputs
	return
}
