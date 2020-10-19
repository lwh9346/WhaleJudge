package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

//StartHTTPServer 启动代码执行服务器
func StartHTTPServer() {
	r := gin.Default()

	//添加服务
	//测评部分
	r.POST("/judge", handleJudgeRequest) //程序评测
	//用户部分
	r.POST("/user/register", handleRegisterRequest)         //注册
	r.POST("/user/login", handleLoginRequest)               //登录
	r.POST("/user/info/get", handleUserInfoRequest)         //用户信息查询
	r.POST("/user/info/edit", handleEditUserInfoRequest)    //用户信息修改
	r.POST("/user/editpassword", handleEditPasswordRequest) //更改密码
	//课程部分
	r.POST("/course/get", handleCourseInfoRequest)        //课程信息查询
	r.POST("/course/addteacher", handleAddTeacherRequest) //课程添加教师
	r.POST("/course/join", handleJoinCourseRequest)       //加入课程
	r.POST("/course/exit", handleExitCourseRequest)       //退出课程
	//问题部分
	r.POST("/question/get", handleQuestionInfoRequest)      //问题信息查询
	r.POST("/question/edit", handleEditQuestionRequest)     //问题信息修改
	r.POST("/question/remove", handleRemoveQuestionRequest) //问题删除

	//r.Run()
	run(r)
}

//GetContainerName 返回一个可用的容器名称（uuid）
func GetContainerName() string {
	return uuid.NewV4().String()
}

func run(r *gin.Engine) {
	srv := &http.Server{Handler: r, Addr: "0.0.0.0:8080"}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
