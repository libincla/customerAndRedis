package message

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func Readpkg(c net.Conn) (mes *Message, err error) {
	var byteslice = make([]byte, 4096)

	n, err := c.Read(byteslice[:4])
	if err == io.EOF {
		fmt.Println("读取客户端失败，原因是", err.Error())
		return nil, err
	}

	//根据pkglen读取消息内容
	var pkglen uint32
	pkglen = binary.BigEndian.Uint32(byteslice[0:4])
	n, err = c.Read(byteslice[:pkglen])
	if n != int(pkglen) || err != nil {
		fmt.Println("读取客户端主体失败，原因是", err.Error())
		return nil, err
	}
	//将byteslice[:pkglen]反序列化成 message.Message结构体
	err = json.Unmarshal(byteslice[:pkglen], &mes)
	if err != nil {
		fmt.Println("json反序列化失败，原因是", err.Error())
		return nil, err
	}
	return mes, nil

}
func WritePkg(conn net.Conn, data []byte) (err error) {
	var pkglen uint32
	pkglen = uint32(len(data))
	var byteslice [4]byte
	binary.BigEndian.PutUint32(byteslice[0:4], pkglen)

	_, err = conn.Write(byteslice[:4])
	if err != nil {
		fmt.Println("写入失败，原因是", err.Error())
		return err
	}
	// 发送消息本体
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("发送消息本体失败，原因是", err.Error())
	}
	fmt.Printf("客户端，发送消息的长度是=%d 内容是%s", len(data), string(data))
	return nil
}

// func WritePkg(conn net.Conn) (mes *Message, err error) {

// }
