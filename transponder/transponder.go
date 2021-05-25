package transponder

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

/*
这是一个 借助中介 server 实现的 p2p
args: 本地要穿透的端口, remote 的转发 port, remote 的 server ip:port
1.	remote: 服务器开启 server 监听端口: 用于接收客户端的穿透请求
2.	local: inner_conn -> remote_ip:server.port, 发送数据{ remote 要开启的端口 trans_port }
	remote: receive,  listen trans_port
3.	client: conn remote_ip:trans_port
	remote: copy trans_port_conn <-> inner_conn
*/

func Server(port int) {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	listener, _ := net.ListenTCP("tcp", addr)
	for {
		conn, _ := listener.AcceptTCP()
		go NewTrans(conn)
	}
}

func NewTrans(conn *net.TCPConn) {
	defer conn.Close()
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 15)
	data := make([]byte, 2)
	conn.Read(data)
	transPort := int(binary.LittleEndian.Uint16(data))
	fmt.Println(transPort)

}

func Client(remoteIp string, remotePort int, localPort int, remoteTransPort int) {
	remoteAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", remoteIp, remotePort))
	client, _ := net.DialTCP("tcp", nil, remoteAddr)
	defer client.Close()
	client.SetKeepAlive(true)
	client.SetKeepAlivePeriod(time.Second * 15)
	data := make([]byte, 2) // unsafe.Sizeof(uint16(1))
	binary.LittleEndian.PutUint16(data, uint16(remoteTransPort))
	fmt.Println(remoteTransPort, data)
	client.Write(data)
}
