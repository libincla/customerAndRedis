package utils

import (
	"customerAndRedis/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

//在这里将读包和写包的方法关联到结构体中
//这样的好处在于，以后想传输了，直接调用这个结构体就有两个方法了
type Transfer struct {
	//分析它有哪些字段
	Conn net.Conn
	Buf  [8096]byte
}

func (t *Transfer) Readpkg() (mes *message.Message, err error) {

	n, err := t.Conn.Read(t.Buf[:4])
	if err == io.EOF {
		fmt.Println("读取客户端失败，原因是", err.Error())
		return nil, err
	}

	//根据pkglen读取消息内容
	var pkglen uint32
	pkglen = binary.BigEndian.Uint32(t.Buf[0:4])
	n, err = t.Conn.Read(t.Buf[:pkglen])
	if n != int(pkglen) || err != nil {
		fmt.Println("读取客户端主体失败，原因是", err.Error())
		return nil, err
	}
	//将byteslice[:pkglen]反序列化成 message.Message结构体
	err = json.Unmarshal(t.Buf[:pkglen], &mes)
	if err != nil {
		fmt.Println("json反序列化失败，原因是", err.Error())
		return nil, err
	}
	return mes, nil

}
func (t *Transfer) WritePkg(data []byte) (err error) {
	var pkglen uint32
	pkglen = uint32(len(data))
	// var byteslice [4]byte
	binary.BigEndian.PutUint32(t.Buf[0:4], pkglen)

	_, err = t.Conn.Write(t.Buf[:4])
	if err != nil {
		fmt.Println("写入失败，原因是", err.Error())
		return err
	}
	// 发送消息本体
	_, err = t.Conn.Write(data)
	if err != nil {
		fmt.Println("发送消息本体失败，原因是", err.Error())
	}
	fmt.Printf("客户端，发送消息的长度是=%d 内容是%s", len(data), string(data))
	return nil
}
