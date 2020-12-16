package funcProcess

import (
	"customerAndRedis/common/message"
	"customerAndRedis/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

//写方法转发消息

func (sp *SmsProcess) SendGroupMes(msg *message.Message) {

	//遍历服务器端的onlineUsers  map[int]*UserProcess
	//将消息转发出去

	//取出mes中的内容
	var smsMes message.SmsMes

	err := json.Unmarshal([]byte(msg.MegData), &smsMes)
	if err != nil {
		fmt.Println(" 反序列化的错误是", err.Error())
		return
	}

	//这里还需要序列化
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化失败，原因是", err.Error())

	}

	for id, up := range userMgr.onlineUsers {

		//这里，要过滤自己，不要再发短信给自己
		if id == smsMes.UserId {
			continue
		}
		sp.SendMesToEachUser(data, up.Conn)
	}

}

func (sp *SmsProcess) SendMesToEachUser(data []byte, conn net.Conn) {
	//传教一个Transfer实例，发送data

	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发群发消息失败，原因是", err.Error())

	}
}
