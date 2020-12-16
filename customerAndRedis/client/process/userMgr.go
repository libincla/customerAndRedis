package process

import (
	"customerAndRedis/client/model"
	"customerAndRedis/common/message"
	"fmt"
)

//这里维护着客户端的用户列表

//客户端要维护的map长这样 map[int]*message.User

var clientOnlineUsers map[int]*message.User = make(map[int]*message.User, 1024)
var CurUser model.CurUser //我们在用户登录成功后，完成对CurUser的初始化

//在客户端显示当前在线的用户

func outputOnlineUser() {
	//遍历clientOnlineUsers
	fmt.Println("当前在线用户列表")
	for id, _ := range clientOnlineUsers {
		fmt.Println("用户id:", id)

	}
}

//编写一个方法，专门来处理返回的NotifyUserStatus
func updateUserStatus(notify *message.NotifyUserStatus) {
	//适当优化
	user, ok := clientOnlineUsers[notify.UserId]
	if !ok { //表示这钱
		user = &message.User{
			UserId: notify.UserId,
			// UserStatus: notify.Status,
		}

	}

	user.UserStatus = notify.Status

	// user = &message.User{
	// 	UserId:     notify.UserId,
	// 	UserStatus: notify.Status,
	// }
	clientOnlineUsers[notify.UserId] = user

	outputOnlineUser()
}
