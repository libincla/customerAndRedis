package process

import (
	"customerAndRedis/common/message"
	"customerAndRedis/server/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
}

//客户端的注册功能

func (up *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	fmt.Println("请输入用户id:")
	fmt.Scanln(&userId)
	fmt.Println("请输入用户密码")
	fmt.Scanln(&userPwd)
	fmt.Println("请输入用户昵称")
	fmt.Scanln(&userName)
	//1.链接到服务器

	conn, err := net.Dial("tcp", "0.0.0.0:8788")
	if err != nil {
		fmt.Println("连接服务器端出现了问题，错误是", err.Error())
	}
	defer conn.Close()

	// 2.通过conn发送消息给服务器端
	var mes message.Message
	mes.MegType = message.RegisterMesType
	//3.创建一个RegisterMes结构体
	var registerMes message.RegisterMes

	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4. 将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("序列化失败，原因是", err.Error())
		return err
	}

	//5. 把data赋值给mes.Data字段
	mes.MegData = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("消息序列化失败，原因是", err.Error())
		return err
	}

	//7.创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送注册消息失败，原因是", err.Error())
		return err
	}
	fromServerMsg, err := tf.Readpkg() // fromServerMsg就是RegisterResMes
	if err != nil {
		fmt.Println("读取服务器端响应失败，原因是", err.Error())
		return err
	}

	//将fromServerMsg的Data部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(fromServerMsg.MegData), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，重新登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
	// if err != nil {
	// 	fmt.Println("反序列化失败，原因是", err.Error())
	// 	return err
	// }

}

func (u *UserProcess) Login(id int, password string) (err error) {
	fmt.Println("请输入用户id:")
	fmt.Scanln(&id)
	fmt.Println("请输入用户密码")
	fmt.Scanln(&password)

	fmt.Printf("你输入的userid=%d pwd=%s\n", id, password)
	// return nil
	conn, err := net.Dial("tcp", "0.0.0.0:8788")
	if err != nil {
		fmt.Println("连接服务器出问题")
	}
	//延迟关闭
	defer conn.Close()

	//创建一个login消息的结构体
	var mes message.Message
	mes.MegType = message.LoginMesType

	// for {}
	loginmeg := message.LoginMes{
		UserId:  id,
		UserPwd: password,
	}

	// mes.MegData = loginmeg

	//将loginmeg序列化
	jsonlogin, err := json.Marshal(loginmeg)

	//再将序列化后的结构给MegData字段
	mes.MegData = string(jsonlogin)
	if err != nil {
		fmt.Println("json序列化失败，原因是", err.Error())
	}

	//最后将整个消息序列化
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("meg json序列化失败，原因是", err.Error())
		return err
	}

	//此时，data就是我们要发送的数据
	//1. 先把data的长度发给服务器， 先获取到data的长度，转换成一个表示长度的byte切片

	err = message.WritePkg(conn, data)
	if err != nil {
		fmt.Println("写入服务器端失败，原因是", err.Error())
	}

	//读取来自服务器端的响应
	var resLoginMsg message.LoginResMes
	tf := &utils.Transfer{
		Conn: conn,
	}

	fromServerMsg, err := tf.Readpkg()
	if err != nil {
		fmt.Println("读取服务器端消息失败，原因是", err.Error())
	}
	fmt.Println(fromServerMsg)
	//反序列化服务器端msg
	err = json.Unmarshal([]byte(fromServerMsg.MegData), &resLoginMsg)
	if err != nil {
		fmt.Println("反序列化失败，原因是", err.Error())
		return err
	}
	fmt.Println(resLoginMsg.Code, resLoginMsg.Error)
	//这里做个对响应吗的判断，如果为200 就表示登录成功，循环显示登录成功的菜单

	if resLoginMsg.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = id
		CurUser.UserStatus = message.UserOnline

		//可以显示当前在线用户的列表，遍历resLoginMsg的UsersId
		fmt.Println("当前在线用户列表如下")
		for _, v := range resLoginMsg.UsersId {

			//如果要求不显示自己在线的话，可以做一个判断，判断id是否为自己
			if v == id {
				continue
			}
			fmt.Println("用户id为", v)

			//完成对客户端的clientOnlineUsers初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			clientOnlineUsers[v] = user
		}

		//这里我们还需要在客户端启动时启动一个协程
		//该协程保持和服务器端的通讯，如果服务器端有数据推送给客户端则接收并显示在客户端的终端
		go serverProcessMes(conn)
		//登录成功
		for {
			ShowMenu()
		}
	}

	//
	// n, err := conn.Read(comeFromData)
	// if err != nil {
	// 	fmt.Println("读取服务器端响应失败", err.Error())
	// }
	// fmt.Println("读取服务器端响应长度为", n)
	// err = json.Unmarshal(comeFromData[:], &fromServerMsg)
	// if err != nil {
	// 	fmt.Println("反序列化服务器端响应失败", err.Error())
	// }
	// fmt.Println(fromServerMsg)

	// msg, err := message.Readpkg(conn)

	// time.Sleep(10 * time.Second)
	// }
	return nil
}
