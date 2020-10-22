package httpserver

import (
	"github.com/gin-gonic/gin"
)

type Homework struct {
	CreatTime int
	Questions []string
	Title     string
}

type CourseInfo struct {
	Title       string
	Description string
	Homeworks   []Homework
	Teachers    []string
	Students    []string
}

func handleCreateCourseRequest(c *gin.Context) {

}

func handleCourseInfoRequest(c *gin.Context) {

}

func handleAddTeacherRequest(c *gin.Context) {

}

func handleExitCourseRequest(c *gin.Context) {
	//课程中最后一名教师退出课程的时候将课程删除
}

func handleJoinCourseRequest(c *gin.Context) {

}
