package port_forward

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var targetAddr *net.TCPAddr

func logP(v ...interface{}) {
	fmt.Printf(fmt.Sprintln(v[0]), v[1:])
	//switch len(v) {
	//case 1:
	//	fmt.Println(v)
	//default:
	//}
}
func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
func logErr(err error, v ...interface{}) {
	if err != nil {
		if typeof(v[1]) == "func(error)" {
			v[1].(func(error))(err)
		} else {
			log.Fatalf(fmt.Sprintln(v[1]), v[2:]...)
		}
		//switch len(err) {
		//case 1:
		//	log.Fatalln(err[0])
		//default:
		//}
	}
}

func Start() {
	var target string
	var port int

	flag.StringVar(&target, "r", "", "the remote server (<host>:<port>)")
	flag.IntVar(&port, "p", 2222, "the proxy port")
	flag.Parse()

	println(strings.Join(os.Args, " | "))
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	logErr(err)

	targetAddr, err = net.ResolveTCPAddr("tcp", target)
	logErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	logErr(err, "Could not start proxy server on %d: %v\n", port, err)

	fmt.Printf("Proxy server running on %d\n", port)
	for {
		conn, err := listener.AcceptTCP()
		logErr(err, "Could not accept client connection", err)

		// 接收到的客户端的连接
		go handleTCPConn(conn)
	}
}

func handleTCPConn(conn *net.TCPConn) {
	defer conn.Close()
	log.Printf("Client '%v' connected!\n", conn.RemoteAddr())

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 15)

	// 建立一个 tcp 连接, 作为一个客户端
	client, err := net.DialTCP("tcp", nil, targetAddr)
	if err != nil {
		log.Println("Could not connect to remote server:", err)
		return
	}

	defer client.Close()
	log.Printf("Connection to server '%v' established!\n", client.RemoteAddr())

	client.SetKeepAlive(true)
	client.SetKeepAlivePeriod(time.Second * 15)

	stop := make(chan bool)

	go func() {
		io.Copy(client, conn)
		println("1")
		stop <- true
	}()

	go func() {
		io.Copy(conn, client)
		println("2")
		stop <- true
	}()

	<-stop
}