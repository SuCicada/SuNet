package main

import (
	"fmt"
	"github.com/txn2/txeh"
)

func WriteHosts(newAddrs []Addr) {
	// 目前的设定是, 只运行于 linux 或 WSL 中
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	for _, addr := range newAddrs {
		hosts.RemoveAddress(addr.host)
		hosts.AddHost(addr.ip.String(), addr.host)
	}
	hfData := hosts.RenderHostsFile()
	fmt.Println(hfData)
	//hosts.Save()
	// or hosts.SaveAs("./test.hosts")
}
