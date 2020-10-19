# Whale Judger

Whale Judger是一个由go编写，基于docker的编程测评平台。

Whale Judger目前支持将golang作为答案的语言来做题，未来将支持cpp、python等语言。

## 开发现状

目前仍在编写后端，前端工作尚未开始。

后端已经完成了基于docker容器的测评功能。

## 部署（目前仍不能部署）

1. 安装docker
2. 安装golang(>=1.15.2)
3. 克隆本仓库
4. 执行`go build`
5. 执行`WhaleJudger init`以初始化
6. 编辑配置文件
7. 执行`WhaleJudger`以运行http服务器