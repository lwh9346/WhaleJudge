package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xujiajun/nutsdb"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

//StartHTTPServer 启动代码执行服务器
func StartHTTPServer() {
	//初始化数据库
	var dberr error
	dbOption := nutsdb.DefaultOptions
	dbOption.Dir = "./db/user"
	userDB, dberr = nutsdb.Open(dbOption)
	if dberr != nil {
		log.Fatal(dberr)
	}
	defer userDB.Close()
	dbOption.Dir = "./db/question"
	questionDB, dberr = nutsdb.Open(dbOption)
	if dberr != nil {
		log.Fatal(dberr)
	}
	defer questionDB.Close()
	dbOption.Dir = "./db/course"
	courseDB, dberr = nutsdb.Open(dbOption)
	if dberr != nil {
		log.Fatal(dberr)
	}
	defer courseDB.Close()

	r := gin.Default()
	r.Use(cors())
	//添加服务
	//演示部分
	r.POST("/demo", handleDemoRequest)
	//测评部分
	r.POST("/judge", handleJudgeRequest) //程序评测
	//用户部分
	r.POST("/user/register", handleRegisterRequest)         //注册
	r.POST("/user/login", handleLoginRequest)               //登录
	r.POST("/user/info/get", handleUserInfoRequest)         //用户信息查询
	r.POST("/user/info/edit", handleEditUserInfoRequest)    //用户信息修改
	r.POST("/user/editpassword", handleEditPasswordRequest) //更改密码
	//课程部分
	r.POST("/course/create", handleCreateCourseRequest)   //创建新课程
	r.POST("/course/get", handleCourseInfoRequest)        //课程信息查询
	r.POST("/course/addteacher", handleAddTeacherRequest) //课程添加教师
	r.POST("/course/join", handleJoinCourseRequest)       //加入课程
	r.POST("/course/exit", handleExitCourseRequest)       //退出课程
	//问题部分
	r.POST("/question/get", handleQuestionInfoRequest)      //问题信息查询
	r.POST("/question/add", handleAddQuestionRequest)       //添加问题
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

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
