package funcProcess

import (
	"encoding/json"

	"customerAndRedis/common/message"
	"customerAndRedis/server/model"
	"customerAndRedis/server/utils"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn

	//这里再添加一个字段，表明该连接属于哪个用户
	UserId int
}

//在这里我们编写通知所有在线用户的方法
//这个userId要通知其他在线用户，他上线了
func (up *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历 onlineUsers,一个个的发送RegisterResMes
	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知（单独写一个方法）
		err := up.NotifyMeOnline(userId)
		if err != nil {
			fmt.Println("通知其他用户失败，原因是", err.Error())
		}
	}
}

func (up *UserProcess) NotifyMeOnline(userId int) (err error) {
	// 组装我们的notifyUserStatusMes
	var mes message.Message
	mes.MegType = message.NotifyUserStatusType

	var notify message.NotifyUserStatus
	notify.UserId = userId
	notify.Status = message.UserOnline

	//将notify序列化
	data, err := json.Marshal(notify)
	if err != nil {
		fmt.Println("notify信息序列化失败，原因是", err.Error())
	}
	mes.MegData = string(data)

	//再将message序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化消息失败，原因是", err.Error())
		return
	}

	//发送，创建tranfser实例
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("transfer 失败，原因是", err.Error())
		return
	}
	return
}

func (up *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//先从mes中取出mes.Data,反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.MegData), &registerMes)
	if err != nil {
		fmt.Println("json反序列化失败，原因是", err.Error())
		return err
	}
	//先声明一个resMes
	var resMes message.Message
	resMes.MegType = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//需要从redis数据库中完成注册
	//使用model.MyUserDao去redis中验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册中发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("序列化失败，原因是", err.Error())
		return err
	}

	//将data 赋值给resMes
	resMes.MegData = string(data)

	//对resMes序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化resmes失败，原因是", err.Error())
		return err
	}

	//发送data，我们将其封装到writePkg函数，因为使用的分层模式，我们先创建Transfer实例，
	// 然后读取
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return

}

//函数serverProcessLogin函数，专门处理用户登录需求
func (u *UserProcess) ServerProcessLogin(msg *message.Message) (err error) {
	// msg, err = message.Readpkg(conn)
	// msg, err = readPkg(conn)
	fmt.Println("msg=", msg)
	//验证内嵌的结构体消息，并返回登陆的响应吗
	var toClientMes message.Message
	var loginMes message.LoginMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(msg.MegData), &loginMes)
	if err != nil {
		fmt.Println("msg中data反序列化失败", err.Error())
		return err
	}

	//需要到redis数据库去验证，直接使用model.UserDao去redis里验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	fmt.Println(user)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
		// 	loginResMes.Code = 500
		// 	loginResMes.Error = "验证失败"
		// } else {
		// 	loginResMes.Code = 200
		// 	loginResMes.Error = ""
	} else {
		loginResMes.Code = 200
		//这里，因为用户登录成功了，就把该登录成功的用户放入userMgr中
		//并将登录成功的用户的UserId赋给u这个UserProcess实例
		u.UserId = loginMes.UserId
		userMgr.AddOnlineUser(u)

		//通知其他在线用户，我上线了
		u.NotifyOthersOnlineUser(u.UserId)

		//将当前在线用户的id，放入loginResMes.UsersId
		//遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}

		fmt.Println(user, "登录成功")
	}
	// if loginMes.UserId != 456 && loginMes.UserPwd != "abc" {
	// 	//验证失败
	// 	loginResMes.Code = 500
	// 	loginResMes.Error = "验证失败"
	// } else {
	// 	loginResMes.Code = 200
	// 	loginResMes.Error = ""
	// }
	//对loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("对响应结构体序列化失败，原因是", err.Error())
		return err
	}
	toClientMes.MegData = string(data)
	toClientMes.MegType = message.LoginResMesType

	data, err = json.Marshal(toClientMes)
	if err != nil {
		fmt.Println("对响应消息序列化失败，原因是", err.Error())
		return err
	}

	//发送data，这里要实例化一个Transfer实例，将它的conn绑定
	tf := &utils.Transfer{
		Conn: u.Conn,
	}

	err = message.WritePkg(tf.Conn, data)
	if err != nil {
		fmt.Println("发送客户端响应出问题，原因是", err.Error())
		return err
	}

	return nil
	// //序列化loginResMes
	// fmt.Println(loginResMes)
	// loginMesByte, err := json.Marshal(loginResMes)
	// if err != nil {
	// 	fmt.Println("响应结构体序列化失败，原因是", err.Error())
	// }
	// toClientMes.MegType = message.LoginResMesType
	// toClientMes.MegData = string(loginMesByte)
	// // 序列化toClientMes
	// toClientData, err := json.Marshal(toClientMes)
	// if err != nil {
	// 	fmt.Println("序列化响应消息失败，原因是", err.Error())
	// }
	// _, err = conn.Write(toClientData)
	// if err != nil {
	// 	fmt.Println("返回客户端失败，原因是", err.Error())
	// } else {
	// 	fmt.Println("发送响应成功")
	// }
	// return
}
