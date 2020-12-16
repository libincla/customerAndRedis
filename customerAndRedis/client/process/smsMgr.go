package process

import (
	"customerAndRedis/common/message"
	"encoding/json"
	"fmt"
)

//专门处理服务器群发过来的消息

func outputGroupMes(mes *message.Message) { //这个地方mes一定是SmsMes
	//显示即可

	//1.反序列化message.Message
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.MegData), &smsMes)
	if err != nil {
		fmt.Println("反序列化失败，原因是", err.Error())
	}
	//显示信息
	info, _ := fmt.Printf("用户id:%d 对大家说%s\n", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	// content = smsMes.
}
