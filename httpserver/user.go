package httpserver

import (
	"encoding/json"
	"log"
	"strings"
	"unicode"

	"github.com/lwh9346/WhaleJudger/database"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

//UserInfo 用户公开信息
//TODO:已完成的题目列表，加入的班级等
type UserInfo struct {
	NickName string `json:"nickname"`
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
		c.JSON(400, gin.H{})
		return
	}
	var ui UserInfo
	ui.NickName = string(database.GetValue(userDB, usernameUserInfoBK, uir.Username))
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

//RegisterResponse 注册用户的结果
type RegisterResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"` //0成功，1失败
}

//handleRegisterRequest 处理注册用户的请求
func handleRegisterRequest(c *gin.Context) {
	var regRequest RegisterRequest
	var regResponse RegisterResponse
	if c.BindJSON(&regRequest) != nil {
		regResponse.Code = 1
		regResponse.Msg = "请求格式不正确"
		c.JSON(400, regResponse)
		return
	}
	//用户名合法性检验
	if len(regRequest.Username) == 0 {
		regResponse.Code = 1
		regResponse.Msg = "没有填写用户名"
		c.JSON(400, regResponse)
		return
	}
	ok := true
	for _, r := range []rune(regRequest.Username) {
		ok = ok && (unicode.IsLetter(r) || unicode.IsNumber(r)) //用户名只能为字母加数字
	}
	if !ok {
		regResponse.Code = 1
		regResponse.Msg = "用户名只能包含字母与数字"
		c.JSON(400, regResponse)
		return
	}
	regRequest.Username = strings.ToLower(regRequest.Username) //不区分大小写
	//查重
	if database.HasKey(userDB, "userpass", regRequest.Username) {
		regResponse.Code = 1
		regResponse.Msg = "用户名已存在"
		c.JSON(400, regResponse)
		return
	}
	//密码合法性检查
	if regRequest.Password == "" {
		regResponse.Code = 1
		regResponse.Msg = "没有填写密码"
		c.JSON(400, regResponse)
		return
	}
	if len([]rune(regRequest.Password)) < 8 {
		regResponse.Code = 1
		regResponse.Msg = "密码太短，至少要8位"
		c.JSON(400, regResponse)
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
	regResponse.Code = 0
	regResponse.Msg = "注册成功"
	c.JSON(200, regResponse)
}

//LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//LoginResponse 登陆请求返回值，返回一个在一天内有效的token用于后续身份验证
type LoginResponse struct {
	Msg   string `json:"msg"`
	Code  int    `json:"code"`
	Token string `json:"token"`
}

func handleLoginRequest(c *gin.Context) {
	var logRequest LoginRequest
	if c.BindJSON(&logRequest) != nil {
		c.JSON(400, gin.H{"msg": "请求格式不正确", "code": 1})
		return
	}
	if !database.HasKey(userDB, usernamePasswordBK, logRequest.Username) {
		c.JSON(401, gin.H{"code": 1, "msg": "用户名不存在"})
		return
	}
	if string(database.GetValue(userDB, usernamePasswordBK, logRequest.Username)) != logRequest.Password {
		c.JSON(401, gin.H{"code": 1, "msg": "密码不正确"})
		log.Println(database.GetValue(userDB, usernamePasswordBK, logRequest.Username))
		return
	}
	token := uuid.NewV4().String()
	database.SetValue(userDB, tokenUsernameBK, token, []byte(logRequest.Username), 3600*24)
	var logResponse LoginResponse
	logResponse.Code = 0
	logResponse.Msg = "登陆成功"
	logResponse.Token = token
	c.JSON(200, logResponse)
}

func handleEditPasswordRequest(c *gin.Context) {

}
