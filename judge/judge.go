package judge

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

//SingleCase 返回值中0表示pass,1表示wa,2表示timeout,3表示error
func SingleCase(containerName, input, stdOutput string, args []string) (output string, statusCode int) {
	args = append([]string{"exec", "-i", containerName}, args...)
	cmd := exec.Command("docker", args...)
	cmd.Stdin = strings.NewReader(input)
	outputChan := make(chan string)
	errChan := make(chan int)
	timeOutChan := time.NewTimer(time.Second).C
	go func() {
		o, e := cmd.Output()
		if e != nil {
			errChan <- 0
			return
		}
		outputChan <- string(o)
	}()
	select {
	case output = <-outputChan:
		outputLines := strings.Split(output, "\n")
		stdOutputLines := strings.Split(stdOutput, "\n")
		if len(stdOutputLines) > len(outputLines) {
			return "WA:你的答案行数小于标准答案", 1
		}
		for k := range stdOutputLines {
			if stdOutputLines[k] != outputLines[k] {
				output = fmt.Sprintf("WA:你的答案第%d行错误\n标准答案：\n%s\n你的答案：\n%s", k+1, stdOutputLines[k], outputLines[k])
				return output, 1
			}
		}
		return "PASS:通过", 0
	case <-errChan:
		return "ERR:程序发生错误", 3
	case <-timeOutChan:
		cmd.Process.Kill()
		return "TO:程序运行超时", 2
	}
}
