package transponder

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
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

const (
	CONN  = 1
	HEART = 2
	CLOSE = 3
	DATA  = 4
	YES   = 0x01
	NEW   = 0x02
)

type event struct {
	signal int16
	value  []byte
}

func Server(port int) {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	listener, _ := net.ListenTCP("tcp", addr)
	log.Println("listen: ", addr)
	for {
		conn, _ := listener.AcceptTCP()
		go NewTrans(conn)
	}
}

func NewTrans(clientConn *net.TCPConn) {
	log.Println("new tans", clientConn.RemoteAddr())
	defer func() {
		log.Println("tans Close", clientConn.RemoteAddr())
		clientConn.Close()
	}()
	clientConn.SetKeepAlive(true)
	clientConn.SetKeepAlivePeriod(time.Second * 3)
	data := make([]byte, 2)
	clientConn.Read(data)
	//clientConn.SetReadDeadline(time.Now().Add(time.Second * 6))
	transPort := int(binary.BigEndian.Uint16(data))
	transAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", transPort))
	// 监听本地的代理端口
	transListener, _ := net.ListenTCP("tcp", transAddr)
	_, err := clientConn.Write([]byte{YES})
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("transListener ", transAddr)

	for {
		//transListener.SetDeadline(time.Now().Add(time.Second * 10))
		go func() {
			time.Sleep(time.Second * 1)
			d := make([]byte, 1)
			for {
				if _, err := clientConn.Read(d); err != nil {
					log.Println(err.Error())
					break
				}
				//fmt.Println("d", d)
			}
			if err := transListener.Close(); err != nil {
				log.Println(err.Error())
			}
			log.Println("transListener.Close()", transListener)
		}()
		// 客户端的连接
		transConn, err := transListener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			break
		}

		// 给内网连接发送信号
		clientConn.Write([]byte{NEW})
		log.Println("clientConn.Write([]byte{NEW})")
		// 接入来自内网的连接
		transClientConn, err := transListener.AcceptTCP()
		if err != nil {
			log.Println(err)
		}
		go newTransConn(transClientConn, transConn,
			func() {
			})
	}
}

/*
clientConn 是内网穿透连接
transConn 是新的 client 请求
*/
func newTransConn(clientConn *net.TCPConn, transConn *net.TCPConn, ready func()) {
	log.Println("newTransConn", transConn.RemoteAddr())
	defer func() {
		if err := transConn.Close(); err != nil {
			log.Println(err.Error())
		}
		if err := clientConn.Close(); err != nil {
			log.Println(err.Error())
		}
		log.Println("transConn Close", transConn.RemoteAddr())
	}()
	clientConn.SetKeepAlive(true)
	clientConn.SetKeepAlivePeriod(time.Second * 5)

	transConn.SetKeepAlive(true)
	transConn.SetKeepAlivePeriod(time.Second * 5)

	stop := make(chan bool)
	go func() {
		// 实验结果表明, tcp连接会在这里阻塞住
		_, err := io.Copy(clientConn, transConn)
		if err != nil {
			log.Println(err)
		}
		stop <- true
	}()

	go func() {
		_, err := io.Copy(transConn, clientConn)
		if err != nil {
			log.Println(err)
		}
		stop <- true
	}()
	ready()
	<-stop
	log.Println("over")
}

//==========================
//type Client struct {
//	localPort int
//}

//func (client *Client) connLocalPool() *net.TCPConn {

//return localClient
//}

func Client(remoteIp string, remotePort int, localPort int, remoteTransPort int) {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	remoteAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", remoteIp, remotePort))
	remoteClient, _ := net.DialTCP("tcp", nil, remoteAddr)
	log.Println(`remoteAddr:`, remoteAddr)
	log.Println(`localPort:`, localPort)
	log.Println(`remoteTransPort:`, remoteTransPort)
	log.Println(`remoteClient:`, remoteClient.LocalAddr(), remoteClient.RemoteAddr())
	defer func() {
		remoteClient.Close()
		log.Println("remoteClient Close")
	}()
	remoteClient.SetKeepAlive(true)
	remoteClient.SetKeepAlivePeriod(time.Second * 5)

	data := make([]byte, 2) // unsafe.Sizeof(uint16(1))
	binary.BigEndian.PutUint16(data, uint16(remoteTransPort))
	//connEvent := event{CONN, data}
	//connEventByte, _ := json.Marshal(connEvent)
	remoteClient.Write(data)
	data = make([]byte, 1)

	// 等待 YES 验证
	remoteClient.Read(data)
	log.Println("data", data)

	if data[0] == YES {
		go func() {
			for {
				time.Sleep(time.Second * 1)
				_, err := remoteClient.Write([]byte{0x00})
				if err != nil {
					log.Println(err.Error())
					return
				} else {
					log.Println("heart")
				}
			}
		}()
		for {
			// 等待 NEW 新连接
			if _, err := remoteClient.Read(data); err != nil {
				//log.Println(err.Error())
			}
			//log.Println("2data", data)
			if data[0] == NEW {
				go func() {
					remoteTransAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", remoteIp, remoteTransPort))
					remoteTransClient, _ := net.DialTCP("tcp", nil, remoteTransAddr)
					localClient := newLocalConn(localPort)
					defer func() {
						localClient.Close()
						remoteTransClient.Close()
						log.Printf("localClient.Close: %s <-> %s\n",
							localClient.LocalAddr(), localClient.RemoteAddr())
						log.Printf("remoteTransClient.Close: %s <-> %s\n",
							remoteTransClient.LocalAddr(), remoteTransClient.RemoteAddr())
					}()
					log.Println("new localClient", localClient.LocalAddr())
					transportData(localClient, remoteTransClient)
				}()
			}
		}
	}
	//go client.connLocalPool()
}
func newLocalConn(localPort int) *net.TCPConn {
	localAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", localPort))
	localClient, _ := net.DialTCP("tcp", nil, localAddr)
	log.Println("localClient", localAddr)
	//defer func() {
	//	localClient.Close()
	//	log.Println("localClient Close")
	//}()
	localClient.SetKeepAlive(true)
	localClient.SetKeepAlivePeriod(time.Second * 5)
	return localClient
}

func transportData(localClient *net.TCPConn, remoteClient *net.TCPConn) {
	stop := make(chan bool)
	go func() {
		if _, err := io.Copy(localClient, remoteClient); err != nil {
			log.Println("io.Copy(localClient, remoteClient)", err.Error())
		}
		//Copy(localClient, remoteClient)
		stop <- true
	}()

	go func() {
		if _, err := io.Copy(remoteClient, localClient); err != nil {
			log.Println("io.Copy(remoteClient, localClient)", err.Error())
		}
		//Copy(remoteClient, localClient)
		stop <- true
	}()
	// 当没有数据流动的时候, 就停止了,
	// 实际上这是不行的, 要一直保活

	<-stop
	log.Println("Client stop")
}
func Copy(reader *net.TCPConn, writer *net.TCPConn) {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf) // 从conn中读取客户端发送的数据内容
		if err != nil {
			if err == io.EOF {
				fmt.Printf("客户端退出 err=%v\n", err)
			} else {
				fmt.Printf("read err=%v\n", err)
			}
			break
		}
		bb := buf[:n]
		fmt.Print(bb)
		writer.Write(bb)
	}
}
