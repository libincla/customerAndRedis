package main

import (
	"customerAndRedis/client/process"
	"fmt"
	"os"
)

var msg string = `
--------欢迎登陆多人聊天系统-----------
1. 登陆聊天系统
2. 注册用户
3. 退出系统

请选择(1-3):
------------
`
var userId int
var userPass string
var userName string
var loop bool = true

func main() {
	var yourchoice int
	for loop {
		fmt.Println(msg)
		fmt.Scanln(&yourchoice)
		switch yourchoice {
		case 1:
			fmt.Println("请登陆")
			// loop = false
			//完成登录
			//1.创建一个UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPass)
		case 2:
			fmt.Println("注册用户")
			up := &process.UserProcess{}
			up.Register(userId, userPass, userName)
		case 3:
			fmt.Println("退出系统")
			loop = false
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")
		}

	}
}
