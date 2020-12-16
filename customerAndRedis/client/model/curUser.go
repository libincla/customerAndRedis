package model

import (
	"customerAndRedis/common/message"
	"net"
)

//当前客户端的User
//因为在客户端，很多地方都会使用到这个CurUser,我们将其声明为全局的
type CurUser struct {
	Conn net.Conn
	message.User
}
