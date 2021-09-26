// https://cloud.tencent.com/developer/article/1075942
// https://developer.51cto.com/art/202101/639962.htm
package lan_group

import (
	"encoding/binary"
	"fmt"
	"github.com/loveleshsharma/gohive"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type IP uint32
type Addr struct {
	ip   net.IP
	host string
}

var addrChan chan Addr
var wg sync.WaitGroup

//线程池大小
var pool_size = 50000
var pool = gohive.NewFixedSizePool(pool_size)

func HttpClient(isWrite bool, scanScope string) {
	var begin = time.Now()

	// 不设置多一些的缓冲区就会报错
	addrChan = make(chan Addr, 100)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本地网络信息:", err)
	}
	var needScanNet []net.IP
	if scanScope == "" {
		for _, a := range addrs {
			needScanNet = append(needScanNet, ScanAddress(a)...)
			//go
			//wg.Wait()
		}
	} else {
		_, ip, _ := net.ParseCIDR(scanScope)
		needScanNet = append(needScanNet, getTable(ip)...)
	}
	fmt.Println("needScanNet: ", len(needScanNet))
	//bar :=progressbar.Default(int64(len(needScanNet)))
	bar := progressbar.NewOptions(
		len(needScanNet),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetDescription("scan..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	for _, ip := range needScanNet {
		wg.Add(1)
		pool.Submit(func() {
			request(ip)
			bar.Add(1)
		})
	}
	wg.Wait()
	close(addrChan)
	fmt.Println("len(addrChan)", len(addrChan))
	if isWrite {
		var newAddrs []Addr
		for arr := range addrChan {
			newAddrs = append(newAddrs, arr)
		}
		WriteHosts(newAddrs)
	}
	var elapseTime = time.Now().Sub(begin)
	fmt.Println("耗时:", elapseTime)
}

func ScanAddress(a net.Addr) []net.IP {
	var needScanNet []net.IP
	if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
		if ipNet.IP.To4() != nil {
			//&& ipNet.IP.String()=="172.27.112.1"
			needScanNet = append(needScanNet, getTable(ipNet)...)
			//wg.Add(1)
			//pool.Submit(func() {
			//})
		}
	}
	//defer wg.Done()
	return needScanNet
}

var HTTPTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second, // 连接超时时间
		KeepAlive: 5 * time.Second, // 保持长连接的时间
	}).DialContext, // 设置连接的参数
	MaxIdleConns:          500,              // 最大空闲连接
	IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
	ExpectContinueTimeout: 30 * time.Second, // 等待服务第一个响应的超时时间
	//MaxIdle/ConnsPerHost:   100,              // 每个host保持的空闲连接数
}
var client = http.Client{
	Transport: HTTPTransport,
	Timeout:   5 * time.Second,
}

func request(ipNet net.IP) {
	url := fmt.Sprintf("http://%s:%d/list", ipNet, port)
	resp, err := client.Post(url, "text/plain",
		strings.NewReader("PLZ"))
	if debugIP != "" && ipNet.String() == debugIP {
		if err == nil {
			fmt.Println("no err")
		} else {
			fmt.Println(err.Error())
		}
	}
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		res := string(body)
		arr := strings.Split(res, ":")
		if arr[0] == "OK" {
			log.Println("見つけた", ipNet, arr[1])
			addrChan <- Addr{ipNet, arr[1]}
		}
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		wg.Done()
	}()
}

// Table 根据IP和mask换算内网IP范围
func getTable(ipNet *net.IPNet) []net.IP {
	ip := ipNet.IP.To4()
	maskArray := ipNet.Mask
	var min, max IP
	var data []net.IP
	var mask uint32
	for i := 0; i < 4; i++ {
		mask = mask<<8 + uint32(maskArray[i])
		b := ip[i] & maskArray[i]
		min = (min << 8) + IP(b)
	}
	max = min + IP(0xffffffff^(mask))
	// max 是广播地址，忽略
	// i & 0x000000ff  == 0 是尾段为0的IP，根据RFC的规定，忽略
	for i := min; i < max; i++ {
		if i&0x000000ff != 0 {
			data = append(data, transIp(i))
		}
	}
	log.Printf("本机ip: %s\n"+
		"\t内网IP范围:%s --- %s\n"+
		"\t共 %d\n",
		ip, transIp(min), transIp(max), len(data))
	return data
}
func transIp(ip IP) net.IP {
	s := make(net.IP, 4)
	//var s net.IP
	binary.BigEndian.PutUint32(s, uint32(ip))
	return s
}
