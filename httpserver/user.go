package httpserver

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/lwh9346/WhaleJudger/database"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

//UserInfo 用户公开信息
//TODO:加入的班级
type UserInfo struct {
	NickName         string   `json:"nickname"`
	CreatedQuestions []string `json:"createdquestions"`
	PassedQuestions  []string `json:"passedquestions"`
	Courses          []string `json:"courses"`
}

//UserInfoRequest 用户的公开信息的请求（目前就一个）
type UserInfoRequest struct {
	Username string `json:"username" binding:"required"`
}

//handleUserInfoRequest 处理用户公开信息的请求
func handleUserInfoRequest(c *gin.Context) {
	var uir UserInfoRequest
	if c.BindJSON(&uir) != nil {
		c.JSON(400, gin.H{})
		return
	}
	if !database.HasKey(userDB, usernameUserInfoBK, uir.Username) {
		c.JSON(404, gin.H{})
		return
	}
	var ui UserInfo
	ui.NickName = string(database.GetValue(userDB, usernameUserInfoBK, uir.Username))
	createdQuestionRawData, ok := database.SAllMembers(userDB, usernameCreatedQuestionsBK, uir.Username)
	if ok {
		ui.CreatedQuestions = make([]string, len(createdQuestionRawData))
		for k := range createdQuestionRawData {
			ui.CreatedQuestions[k] = string(createdQuestionRawData[k])
		}
	}
	passedQuestionRawData, ok := database.SAllMembers(userDB, usernamePassedQuestionsBK, uir.Username)
	if ok {
		ui.PassedQuestions = make([]string, 0, len(passedQuestionRawData))
		for k := range passedQuestionRawData {
			v := string(passedQuestionRawData[k])
			if database.HasKey(questionDB, questionDescriptionBK, v) {
				ui.PassedQuestions = append(ui.PassedQuestions, v)
			} else {
				database.SRemove(userDB, usernamePassedQuestionsBK, uir.Username, passedQuestionRawData[k])
			}
		}
	}
	courses, ok := database.SAllMembers(userDB, usernameCourseNamesBK, uir.Username)
	if ok {
		ui.Courses = make([]string, len(courses))
		for k := range courses {
			ui.Courses[k] = string(courses[k])
		}
	}
	c.JSON(200, ui)
}

//EditUserInfoRequest 修改昵称的请求
type EditUserInfoRequest struct {
	Token       string `json:"token" binding:"required"`
	NewNickname string `json:"nickname" binding:"required"`
}

func handleEditUserInfoRequest(c *gin.Context) {
	var euir EditUserInfoRequest
	if c.BindJSON(&euir) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, euir.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, euir.Token))
	database.SetValue(userDB, usernameUserInfoBK, username, []byte(euir.NewNickname), 0)
	c.JSON(200, gin.H{"code": 0, "msg": "昵称修改成功"})
}

//RegisterRequest 注册用户的请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

//handleRegisterRequest 处理注册用户的请求
func handleRegisterRequest(c *gin.Context) {
	var regRequest RegisterRequest
	if c.BindJSON(&regRequest) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	//用户名合法性检验
	if len(regRequest.Username) > 64 {
		c.JSON(400, gin.H{"code": 1, "msg": "用户名过长"})
		return
	}
	ok := true
	for _, r := range []rune(regRequest.Username) {
		ok = ok && (unicode.IsLetter(r) || unicode.IsNumber(r)) //用户名只能为字母加数字
	}
	if !ok {
		c.JSON(400, gin.H{"code": 1, "msg": "用户名只能包含字母与数字"})
		return
	}
	regRequest.Username = strings.ToLower(regRequest.Username) //不区分大小写
	//查重
	if database.HasKey(userDB, "userpass", regRequest.Username) {
		c.JSON(400, gin.H{"code": 1, "msg": "用户名已存在"})
		return
	}
	//密码合法性检查
	if regRequest.Password == "" {
		c.JSON(400, gin.H{"code": 1, "msg": "没有密码"})
		return
	}
	if len([]rune(regRequest.Password)) < 8 {
		c.JSON(400, gin.H{"code": 1, "msg": "密码太短，至少要8位"})
		return
	}
	//昵称检查
	if regRequest.Nickname == "" {
		regRequest.Nickname = regRequest.Username
	}
	//写入数据
	database.SetValue(userDB, usernamePasswordBK, regRequest.Username, []byte(regRequest.Password), 0)
	var ui UserInfo
	ui.NickName = regRequest.Nickname
	uib, _ := json.Marshal(ui)
	database.SetValue(userDB, usernameUserInfoBK, regRequest.Username, uib, 0)
	//返回数据
	c.JSON(200, gin.H{"code": 0, "msg": "注册成功"})
}

//LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func handleLoginRequest(c *gin.Context) {
	var logRequest LoginRequest
	if c.BindJSON(&logRequest) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, usernamePasswordBK, logRequest.Username) {
		c.JSON(401, gin.H{"code": 1, "msg": "用户名不存在"})
		return
	}
	if string(database.GetValue(userDB, usernamePasswordBK, logRequest.Username)) != logRequest.Password {
		c.JSON(401, gin.H{"code": 1, "msg": "密码不正确"})
		return
	}
	token := uuid.NewV4().String()
	database.SetValue(userDB, tokenUsernameBK, token, []byte(logRequest.Username), 3600*24)
	c.JSON(200, gin.H{"code": 0, "msg": "登陆成功", "token": token})
}

//EditPasswordRequest 修改密码的请求
type EditPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"oldpassword" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required"`
}

func handleEditPasswordRequest(c *gin.Context) {
	var epr EditPasswordRequest
	if c.BindJSON(&epr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, usernamePasswordBK, epr.Username) {
		c.JSON(401, gin.H{"code": 1, "msg": "用户名不存在"})
		return
	}
	if string(database.GetValue(userDB, usernamePasswordBK, epr.Username)) != epr.OldPassword {
		c.JSON(401, gin.H{"code": 1, "msg": "密码不正确"})
		return
	}
	if len([]rune(epr.NewPassword)) < 8 {
		c.JSON(400, gin.H{"code": 1, "msg": "密码太短，至少要8位"})
		return
	}
	database.SetValue(userDB, usernamePasswordBK, epr.Username, []byte(epr.NewPassword), 0)
	c.JSON(200, gin.H{"code": 0, "msg": "密码修改成功"})
}
