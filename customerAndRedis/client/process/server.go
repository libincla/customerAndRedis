package process

import (
	"customerAndRedis/common/message"
	"customerAndRedis/server/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面...
func ShowMenu() {
	fmt.Println("----登录成功-----")
	fmt.Println("----1.显示用户列表-----")
	fmt.Println("----2.发送消息-----")
	fmt.Println("----3.信息列表-----")
	fmt.Println("----4.退出系统-----")
	fmt.Println("----请选择(1-4)----")

	var key int
	var content string //要对大家说的话
	fmt.Scanln(&key)
	//因为，我们总会使用到SmsProcess实例，因此我们将其定义在switch外面，
	var smsProcess *SmsProcess = &SmsProcess{}
	switch key {
	case 1:
		outputOnlineUser()
	case 2:
		fmt.Println("你想对大家说什么? ")
		fmt.Scanln(&content)
		smsProcess.SendGroupMes(content)
	case 3:
	case 4:
		fmt.Println("选择退出了系统")
		os.Exit(0)
	default:
		fmt.Println("输入有误")
	}
}

//和服务器端通信

func serverProcessMes(conn net.Conn) {
	//创建一个transfer实例，不停的读取服务器端发过来的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器端发送过来的消息")
		mes, err := tf.Readpkg()
		if err != nil {
			fmt.Println("读取pkg失败，原因是", err.Error())
		}
		// //如果读取到消息，做下一步处理
		// fmt.Printf("mes=%v\n", mes)

		switch mes.MegType {
		case message.NotifyUserStatusType: //有人状态发生了变化
			// 0. 客户端初始化好了 这个map
			// 1. 取出notify消息中的UserId和Status
			// 2. 把这个用户的消息，状态保存到客户端的map中
			//处理
			var notify message.NotifyUserStatus
			json.Unmarshal([]byte(mes.MegData), &notify)
			updateUserStatus(&notify)

		case message.SmsMesType:
			//有人群发消息了
			outputGroupMes(mes)
		default:
			fmt.Println("服务器端返回了一个不能识别的消息类型")
		}

	}
}
