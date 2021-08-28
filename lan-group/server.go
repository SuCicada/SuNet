package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func HttpServer(port int) {
	http.HandleFunc("/list",
		func(w http.ResponseWriter, r *http.Request) {
			s, _ := ioutil.ReadAll(r.Body)
			log.Println(r)
			log.Println(s)
			if bytes.Compare(s, PLZ) == 0 {
				log.Println("get", r.Header["User-Agent"], r.Method, r.RemoteAddr)
				res := fmt.Sprintf("OK:%s", getSelfHost())
				if _, err := fmt.Fprint(w, res); err != nil {
					log.Println(err.Error())
				}
			}
		})
	log.Println("now listen:", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panicln(err.Error())
	}
	fmt.Println("ok")
}

func SelfInfo() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本地网络信息:", err)
	}
	for index, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(index, "IP:", ipNet.IP)
				fmt.Println("   子网掩码:", ipNet.Mask)
				it, err := net.InterfaceByIndex(index)
				if err == nil {
					fmt.Println("Mac地址:", it.HardwareAddr)
				}
			}
		}
	}
}
