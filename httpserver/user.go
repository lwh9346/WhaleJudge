package httpserver

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/lwh9346/WhaleJudger/database"

	"github.com/gin-gonic/gin"
)

//UserInfo 用户公开信息
type UserInfo struct {
	NickName string `json:"nickname"`
}

func handleUserInfoRequest(c *gin.Context) {

}

func handleEditUserInfoRequest(c *gin.Context) {

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

func handleLoginRequest(c *gin.Context) {

}

func handleEditPasswordRequest(c *gin.Context) {

}
