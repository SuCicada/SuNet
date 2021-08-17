// https://cloud.tencent.com/developer/article/1075942
package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func HttpClient() {
	// 单网卡模式
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本地网络信息:", err)
	}
	for i, a := range addrs {
		ScanAddress(i, a)
	}
}
func ScanAddress(index int,a net.Addr) {
	if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
		if ipNet.IP.To4() != nil {
			fmt.Println("IP:", ipNet.IP)
			fmt.Println("子网掩码:", ipNet.Mask)
			it, err := net.InterfaceByIndex(index)
			if err == nil {
				fmt.Println("Mac地址:", it.HardwareAddr)
			}

			//getTable(ipNet)
			for i, ipInt := range getTable(ipNet) {
				ip := strIp(ipInt)
				go func() {
					res := request(ip)
					if res != "" {
						fmt.Println(i, ip, res)
					}
				}()
			}
		}
	}
}

func request(ipNet net.IP) string {
	url := fmt.Sprintf("http://%s:%d/list", ipNet, port)
	resp, err := http.Post(url, "text/plain",
		strings.NewReader("PLZ"))
	if err == nil {

		body, _ := ioutil.ReadAll(resp.Body)
		res := string(body)
		arr := strings.Split(res, ":")
		if arr[0] == "OK" {
			return arr[1]
		}
	}
	return ""
}

type IP uint32

// Table 根据IP和mask换算内网IP范围
func getTable(ipNet *net.IPNet) []IP {
	ip := ipNet.IP.To4()
	maskArray := ipNet.Mask
	log.Println("本机ip:", ip)
	var min, max IP
	var data []IP
	var mask uint32
	for i := 0; i < 4; i++ {
		mask = mask<<8 + uint32(maskArray[i])
		b := ip[i] & maskArray[i]
		min = (min << 8) + IP(b)
	}
	max = min + IP(0xffffffff^(mask))
	log.Printf("内网IP范围:%s --- %s\n", strIp(min), strIp(max))
	// max 是广播地址，忽略
	// i & 0x000000ff  == 0 是尾段为0的IP，根据RFC的规定，忽略
	for i := min; i < max; i++ {
		if i&0x000000ff == 0 {
			continue
		}
		data = append(data, i)
	}
	return data
}
func strIp(ip IP) net.IP {
	s := make(net.IP, 4)
	//var s net.IP
	binary.BigEndian.PutUint32(s, uint32(ip))
	return s
}
