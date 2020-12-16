package main

import (
	"customerAndRedis/common/message"
	"customerAndRedis/server/funcProcess"
	"customerAndRedis/server/utils"
	"fmt"
	"io"
	"net"
)

//创建一个processor的结构体
type ProcessorInfo struct {
	Conn net.Conn
}

func (p *ProcessorInfo) mainProcess() (err error) {
	for {
		//创建一个Tranfser实例完成读报任务
		tf := &utils.Transfer{
			Conn: p.Conn,
		}
		msg, err := tf.Readpkg()
		// msg, err := readPkg(c)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出")
				return err
			} else {
				fmt.Println("read pkg err", err.Error())
				return err
			}

		}
		err = p.serverProcessMes(msg)
		if err != nil {
			fmt.Println("处理客户端服务出现问题，原因是", err.Error())
			break
		}
	}
	return
}

//根据客户端发送消息种类不同，从而决定调用哪个函数来处理
func (p *ProcessorInfo) serverProcessMes(msg *message.Message) (err error) {

	//看看是否能从客户端发送过来的群发消息
	fmt.Println("mes=", msg)
	switch msg.MegType {
	case message.LoginMesType:

		//处理登陆逻辑
		//创建一个UserProcess实例
		up := &funcProcess.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(msg)
	case message.LoginResMesType:
		fmt.Println("...")
	case message.RegisterMesType:
		//处理注册的逻辑
		//创建一个UserProcess实例
		up := &funcProcess.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRegister(msg)

	case message.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊的消息
		smsProcess := &funcProcess.SmsProcess{}
		smsProcess.SendGroupMes(msg)

	default:
		fmt.Println("pass")
	}
	return err
}
