package httpserver

import (
	"github.com/xujiajun/nutsdb"
)

//注：写在前面的是key，后面的是value

//下面是用户数据部分
var userDB *nutsdb.DB //存放用户数据的数据库

const (
	usernamePasswordBK         = "userpass"           //用户名密码数据库
	tokenUsernameBK            = "tokenuser"          //用于存储token的数据库
	usernamePassedQuestionsBK  = "userpassquestion"   //记录某个用户已经PASS的题目的库，应该用Set结构
	usernameCreatedQuestionsBK = "usercreatequestion" //记录某个用户创建的题目的库，应该用Set结构
	usernameCourseNamesBK      = "usercourse"         //记录某个用户参加的课程的库，应该用Set结构
	usernameUserInfoBK         = "userinfo"           //用户数据，如昵称等的库
)

//下面是题目数据部分
var questionDB *nutsdb.DB //存放题目数据的数据库

const (
	questionDescriptionBK = "questiondescription" //存储题目描述信息以及例题的数据库
	questionCasesBK       = "questioncases"       //存储题目的cases的数据库
)

//下面是课程数据部分
var courseDB *nutsdb.DB //存放课程数据的数据库

/*
var courseDescriptionDB *nutsdb.DB //课程简介的数据库
var courseTeachersDB *nutsdb.DB    //记录课程教师的数据库，一般来说教师较少，不用Set结构
var courseStudentsDB *nutsdb.DB    //记录参与课程的学生的数据库，使用Set结构
*/
