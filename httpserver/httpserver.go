package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

//JudgeRequest 一个测评请求
type JudgeRequest struct {
	//Language 程序所使用的语言
	Language string `json:"language" binding:"required"`
	//SourceCode 程序源代码
	SourceCode string `json:"sourcecode" binding:"required"`
	//QuestionName 所回答的问题的名称（ID）
	QuestionName string `json:"question" binding:"required"`
}

//JudgeResponse 测评结果
type JudgeResponse struct {
	//Msg 包含了每个case的信息，可以是错误信息也可以是其他信息
	Msg []string `json:"msg"`
	//Code 状态码，详见judge.go
	Code int `json:"code"`
}

//StartHTTPServer 启动代码执行服务器
func StartHTTPServer() {
	r := gin.Default()

	//添加服务
	r.POST("/judge", handleJudgeRequest)

	//r.Run()
	run(r)
}

//getInputAndOutputByQuestionName 目前是这样，测试用
func getInputAndOutputByQuestionName(questionName string) (input []string, output []string) {
	input = []string{"hello\n"}
	output = []string{"hello"}
	return
}

func getContainerName() string {
	return "codingTest"
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
