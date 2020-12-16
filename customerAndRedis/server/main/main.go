package main

import (
	"customerAndRedis/server/model"
	"fmt"
	"net"
	"time"
)

// func readPkg(c net.Conn) (mes *message.Message, err error) {
// 	var byteslice = make([]byte, 4096)
// 	fmt.Println("读取客户端")
// 	n, err := c.Read(byteslice[:4])
// 	if err == io.EOF {
// 		fmt.Println("读取客户端失败，原因是", err.Error())
// 		return
// 	}

// 	//根据pkglen读取消息内容
// 	var pkglen uint32
// 	pkglen = binary.BigEndian.Uint32(byteslice[0:4])
// 	n, err = c.Read(byteslice[:pkglen])
// 	if n != int(pkglen) || err != nil {
// 		fmt.Println("读取客户端主体失败，原因是", err.Error())
// 		return
// 	}
// 	//将byteslice[:pkglen]反序列化成 message.Message结构体
// 	err = json.Unmarshal(byteslice[:pkglen], &mes)
// 	if err != nil {
// 		fmt.Println("json反序列化失败，原因是", err.Error())
// 		return
// 	}
// 	return

// }

// }

func goroutineProcess(conn net.Conn) (err error) {
	defer conn.Close()
	processor := &ProcessorInfo{
		Conn: conn,
	}
	err = processor.mainProcess()
	if err != nil {
		return err
	}
	return nil
}

//这里我们编写一个函数，专门完成对UserDao实例的初始化任务

func initUserDao() {
	//这里涉及到初始化的顺序问题，initUserDao一定要在initPool后面执行
	model.MyUserDao = model.NewUserDao(pool)
	// return model.NewUserDao(pool)
}

func main() {
	//当服务器启动时，就开始初始化连接池
	initPool("192.168.0.110:6379", 10, 0, 300*time.Second)
	initUserDao()
	fmt.Println("new server")
	listener, err := net.Listen("tcp", "0.0.0.0:8788")
	if err != nil {
		fmt.Println("监听失败，原因是", err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("conn err", err.Error())
		}
		// fmt.Println(conn)

		go goroutineProcess(conn)
	}

}
