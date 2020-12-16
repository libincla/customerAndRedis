package model

import (
	"customerAndRedis/common/message"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

//我们在服务器启动后，就初始化一个userDao实例，把它做成全局的变量，在需要和redis操作时，就直接使用即可
var (
	MyUserDao *UserDao
)

//定义一个UserDao结构体，它对User结构体的各种操作

type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式，创建一个UserDao的实例
func NewUserDao(pool *redis.Pool) (userdao *UserDao) {
	userdao = &UserDao{
		pool: pool,
	}
	return userdao
}

//要有个方法，可以根据ID来返回一个User的实例+error

func (userdao *UserDao) GetUserById(conn redis.Conn, id int) (user *message.User, err error) {
	//给定id去查询用户
	_, err = conn.Do("AUTH", "yourpass")
	if err != nil {
		fmt.Println("认证失败，原因是", err.Error())
	}
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		//表示出问题
		if err == redis.ErrNil {
			//表示在hash结构中没有找到这个对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &message.User{}

	//这里要把res反序列化成User实例，才能取出User中的各种属性
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json unmarshal err", err.Error())
		return
	}

	return

}

//完成登录的效验Login
//1.Login完成对用户的验证
//2.如果用户的id和pwd都正确，返回user实例，否则返回对应的错误信息
func (userdao *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {

	//先从连接池中取出一个连接
	conn := userdao.pool.Get()
	defer conn.Close()

	//取出反序列化好的用户实例出来
	user, err = userdao.GetUserById(conn, userId)
	if err != nil {
		return
	}

	//下面验证密码是否有效
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (userdao *UserDao) Register(user *message.User) (err error) {

	//先从userDao的连接池中取出一个连接
	conn := userdao.pool.Get()
	defer conn.Close()

	_, err = userdao.GetUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	// 说明id在redis中还没有，可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("序列化user时失败，原因是", err.Error())
		return
	}

	//入库
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户信息错误，原因是", err)
		return
	}
	return
}
