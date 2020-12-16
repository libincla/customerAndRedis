package process

import (
	"customerAndRedis/common/message"
	"customerAndRedis/server/utils"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊的消息

func (sp *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个message
	var mes message.Message
	mes.MegType = message.SmsMesType

	//2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("序列化smsmes失败", err.Error())
		return err
	}
	mes.MegData = string(data)

	//对mes再次序列化
	json.Marshal(mes)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化mes失败", err.Error())
		return err
	}

	//将序列化后的mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("send message error", err.Error())
		return err
	}
	return

}
