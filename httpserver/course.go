package httpserver

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lwh9346/WhaleJudger/database"
)

//Homework 记录了某次作业的信息
type Homework struct {
	CreatTime int      `json:"createtime"`
	Questions []string `json:"questions"`
	Title     string   `json:"title"`
}

//CourseInfo 记录了某个课程的所有信息
type CourseInfo struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Homeworks   []Homework `json:"homeworks" binding:"required"`
	Teachers    []string   `json:"teachers"`
	Students    []string   `json:"Students"`
}

//CreateCourseRequest 创建课程的请求
type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Token       string `json:"token" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func handleCreateCourseRequest(c *gin.Context) {
	courseInfoLock.Lock()
	defer courseInfoLock.Unlock()
	var ccr CreateCourseRequest
	if c.BindJSON(&ccr) != nil {
		c.JSON(400, gin.H{"msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, ccr.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	if len(ccr.Title) > 64 {
		c.JSON(400, gin.H{"code": 1, "msg": "标题过长，最长只能为64个字符"})
		return
	}
	if database.HasKey(courseDB, courseInfoBK, ccr.Title) {
		c.JSON(400, gin.H{"code": 1, "msg": "该课程名已存在"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, ccr.Token))
	ci := CourseInfo{Title: ccr.Title, Description: ccr.Description, Teachers: []string{username}}
	data, _ := json.Marshal(ci)
	database.SetValue(courseDB, courseInfoBK, ccr.Title, data, 0)
	database.SAdd(userDB, usernameCourseNamesBK, username, []byte(ccr.Title))
	c.JSON(200, gin.H{"code": 0, "msg": "课程创建成功"})
}

//CourseInfoRequest 获取课程信息的请求
type CourseInfoRequest struct {
	CourseName string `json:"coursename" binding:"required"`
}

func handleCourseInfoRequest(c *gin.Context) {
	var cir CourseInfoRequest
	if c.BindJSON(&cir) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(courseDB, courseInfoBK, cir.CourseName) {
		c.JSON(404, gin.H{"code": 1, "msg": "找不到该课程"})
		return
	}
	data := database.GetValue(courseDB, courseInfoBK, cir.CourseName)
	var ci CourseInfo
	json.Unmarshal(data, &ci)
	c.JSON(200, ci)
}

//AddTeacherRequest 添加教师的请求
type AddTeacherRequest struct {
	CourseName string `json:"coursename" binding:"required"`
	Token      string `json:"token" binding:"required"`
	UserToAdd  string `json:"usertoadd" binding:"required"`
}

func handleAddTeacherRequest(c *gin.Context) {
	courseInfoLock.Lock()
	defer courseInfoLock.Unlock()
	var atr AddTeacherRequest
	if c.BindJSON(&atr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, atr.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	if !database.HasKey(courseDB, courseInfoBK, atr.CourseName) {
		c.JSON(404, gin.H{"code": 1, "msg": "找不到该课程"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, atr.Token))
	ciData := database.GetValue(courseDB, courseInfoBK, atr.CourseName)
	var ci CourseInfo
	json.Unmarshal(ciData, &ci)
	permitted, _ := isTeacherOfCourse(username, ci)
	if !permitted {
		c.JSON(401, gin.H{"code": 1, "msg": "你没有修改该课程的权限"})
		return
	}
	//至此鉴权完毕
	if !database.HasKey(userDB, usernameUserInfoBK, atr.UserToAdd) {
		c.JSON(400, gin.H{"code": 1, "msg": "你要添加的用户不存在"})
		return
	}
	if is, _ := isTeacherOfCourse(atr.UserToAdd, ci); is {
		c.JSON(400, gin.H{"code": 1, "msg": "你要添加的用户已经是课程教师了"})
		return
	}
	if is, k := isStudentOfCourse(atr.UserToAdd, ci); is {
		ci.Students = append(ci.Students[:k], ci.Students[k+1:]...)
	} else {
		database.SAdd(userDB, usernameCourseNamesBK, atr.UserToAdd, []byte(ci.Title))
	}
	ci.Teachers = append(ci.Teachers, atr.UserToAdd)
	ciData, _ = json.Marshal(ci)
	database.SetValue(courseDB, courseInfoBK, ci.Title, ciData, 0)
	c.JSON(200, gin.H{"code": 0, "msg": "添加成功"})
	return
}

//ExitCourseRequest 退出课程的请求
type ExitCourseRequest struct {
	Token      string `json:"token" binding:"required"`
	CourseName string `json:"coursename" binding:"required"`
}

func handleExitCourseRequest(c *gin.Context) {
	courseInfoLock.Lock()
	defer courseInfoLock.Unlock()
	var ecr ExitCourseRequest
	if c.BindJSON(&ecr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, ecr.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	if !database.HasKey(courseDB, courseInfoBK, ecr.CourseName) {
		c.JSON(404, gin.H{"code": 1, "msg": "找不到该课程"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, ecr.Token))
	ciData := database.GetValue(courseDB, courseInfoBK, ecr.CourseName)
	var ci CourseInfo
	json.Unmarshal(ciData, &ci)
	iss, indexS := isStudentOfCourse(username, ci)
	ist, indexT := isTeacherOfCourse(username, ci)
	if !(iss || ist) {
		c.JSON(400, gin.H{"code": 1, "msg": "你尚未加入该课程"})
		return
	}
	if iss {
		ci.Students = append(ci.Students[:indexS], ci.Students[indexS+1:]...)
	}
	if ist {
		if len(ci.Teachers) == 1 {
			//课程中最后一名教师退出课程的时候将课程删除
			removeCourse(ci)
			c.JSON(200, gin.H{"code": 0, "msg": "已退出该课程，由于你是该课程的唯一一名老师，该课程已删除"})
			return
		}
		ci.Teachers = append(ci.Teachers[:indexT], ci.Teachers[indexT+1:]...)
	}
	ciData, _ = json.Marshal(ci)
	database.SetValue(courseDB, courseInfoBK, ci.Title, ciData, 0)
	database.SRemove(userDB, usernameCourseNamesBK, username, []byte(ci.Title))
	c.JSON(200, gin.H{"code": 0, "msg": "已退出该课程"})
}

//JoinCourseRequest 加入课程的请求
type JoinCourseRequest struct {
	Token      string `json:"token" binding:"required"`
	CourseName string `json:"coursename" binding:"required"`
}

func handleJoinCourseRequest(c *gin.Context) {
	courseInfoLock.Lock()
	defer courseInfoLock.Unlock()
	var jcr JoinCourseRequest
	if c.BindJSON(&jcr) != nil {
		c.JSON(400, gin.H{"code": 1, "msg": "请求格式不正确"})
		return
	}
	if !database.HasKey(userDB, tokenUsernameBK, jcr.Token) {
		c.JSON(401, gin.H{"code": 1, "msg": "登陆失效，请重新登陆"})
		return
	}
	if !database.HasKey(courseDB, courseInfoBK, jcr.CourseName) {
		c.JSON(404, gin.H{"code": 1, "msg": "找不到该课程"})
		return
	}
	username := string(database.GetValue(userDB, tokenUsernameBK, jcr.Token))
	ciData := database.GetValue(courseDB, courseInfoBK, jcr.CourseName)
	var ci CourseInfo
	json.Unmarshal(ciData, &ci)
	iss, _ := isStudentOfCourse(username, ci)
	ist, _ := isTeacherOfCourse(username, ci)
	if iss || ist {
		c.JSON(400, gin.H{"code": 1, "msg": "你已加入该课程了"})
		return
	}
	ci.Students = append(ci.Students, username)
	ciData, _ = json.Marshal(ci)
	database.SetValue(courseDB, courseInfoBK, ci.Title, ciData, 0)
	database.SAdd(userDB, usernameCourseNamesBK, username, []byte(ci.Title))
	c.JSON(200, gin.H{"code": 0, "msg": "成功加入课程"})
}

func isTeacherOfCourse(username string, ci CourseInfo) (is bool, index int) {
	for k, v := range ci.Teachers {
		if v == username {
			return true, k
		}
	}
	return false, -1
}

func isStudentOfCourse(username string, ci CourseInfo) (is bool, index int) {
	for k, v := range ci.Students {
		if v == username {
			return true, k
		}
	}
	return false, -1
}

func removeCourse(ci CourseInfo) {
	for _, v := range ci.Students {
		database.SRemove(userDB, usernameCourseNamesBK, v, []byte(ci.Title))
	}
	for _, v := range ci.Teachers {
		database.SRemove(userDB, usernameCourseNamesBK, v, []byte(ci.Title))
	}
	database.RemoveKey(courseDB, courseInfoBK, ci.Title)
}
