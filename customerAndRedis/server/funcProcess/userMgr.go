package funcProcess

import (
	"fmt"
)

var (
	userMgr *UserMgr
)

//因为UserMgr实例在服务器端有且只有一个
//因为在很多地方都会使用到，因此，我们将其定义为一个全局的变量
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化工作

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUsers的添加
func (ug *UserMgr) AddOnlineUser(up *UserProcess) {
	ug.onlineUsers[up.UserId] = up
}

//删除
func (ug *UserMgr) DelOnlineUser(userid int) {
	delete(ug.onlineUsers, userid)
}

//返回当前所有在线的用户
func (ug *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return ug.onlineUsers
}

//根据id值返回对应的UserProcess，用于点对点的通信
func (ug *UserMgr) GetOnlineUserById(userid int) (up *UserProcess, err error) {

	//从map中取出一个值，带检测的方式
	up, ok := ug.onlineUsers[userid]
	if !ok { //表示获取的用户不在线
		err = fmt.Errorf("用户%d不在线", userid)
		return nil, err
	}
	return up, nil

}
