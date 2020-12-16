package message

const (
	LoginMesType         = "LoginMes"
	LoginResMesType      = "LoginResMes"
	RegisterMesType      = "RegisterMes"
	RegisterResMesType   = "RegisterResMes"
	NotifyUserStatusType = "NotifyUserStatusMes"
	SmsMesType           = "SmsMes"
)

//这里我们定义几个用户状态的常量
const (
	UserOnline     = 0
	UserOffline    = 1
	UserBusyStatus = 2
)

type Message struct {
	MegType string `json: "mesType"`
	MegData string `json: "mesdata`
}

type LoginMes struct {
	UserId  int    `json: "userId"`
	UserPwd string `json: "userPwd"`
}
type LoginResMes struct {
	Code    int    //返回状态码 500 表示用户未注册，200表示登陆成功
	UsersId []int  //增加字段，保存用户id的切片, 将返回在线用户id
	Error   string //返回的错误信息
}

type RegisterMes struct {
	User User `json: "user"` //这个类型就是User结构体
}
type RegisterResMes struct {
	Code  int    `json: "code"`  //返回状态码 400 表示用户已经存在 200表示注册成功
	Error string `json: "error"` //返回的错误信息
}

//定义一个新的消息类型
//为了配合服务器端上线推送用户状态变化的消息
type NotifyUserStatus struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

//增加一个SmsMes //发送的结构体
type SmsMes struct {
	Content string `json:"content"` //内容
	User           //嵌套了一个User的匿名结构体
}
