package judge

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	//Pass 程序通过
	Pass = iota
	//ServerError 服务器错误
	ServerError
	//WrongAnswer 答案错误
	WrongAnswer
	//TimeOut 程序超时(1s)
	TimeOut
	//ProgramError 程序执行错误
	ProgramError
	//CompileError 编译错误
	CompileError
)

//SingleCase 返回值中0表示pass,1表示wa,2表示timeout,3表示error
func SingleCase(containerName, input, stdOutput string, args []string) (output string, statusCode int) {
	args = append([]string{"exec", "-i", containerName}, args...)
	cmd := exec.Command("docker", args...)
	cmd.Stdin = strings.NewReader(input)
	outputChan := make(chan string)
	errChan := make(chan error)
	timeOutChan := time.NewTimer(time.Second).C
	go func() {
		o, e := cmd.Output()
		if e != nil {
			errChan <- e
			return
		}
		outputChan <- string(o)
	}()
	select {
	case output = <-outputChan:
		outputLines := strings.Split(output, "\n")
		stdOutputLines := strings.Split(stdOutput, "\n")
		if len(stdOutputLines) > len(outputLines) {
			return "WA:你的答案行数小于标准答案", WrongAnswer
		}
		for k := range stdOutputLines {
			if stdOutputLines[k] != outputLines[k] {
				output = fmt.Sprintf("WA:你的答案从第%d行开始错误\n标准答案：\n%s\n你的答案：\n%s", k+1, stdOutputLines[k], outputLines[k])
				return output, WrongAnswer
			}
		}
		return "PASS:通过", Pass
	case err := <-errChan:
		return fmt.Sprintf("ERR:程序发生错误\n错误信息：\n%v", err), ProgramError
	case <-timeOutChan:
		cmd.Process.Kill()
		return "TO:程序运行超时", TimeOut
	}
}
