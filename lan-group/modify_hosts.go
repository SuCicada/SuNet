package lan_group

import (
	"fmt"
	"github.com/txn2/txeh"
)

func WriteHosts(newAddrs []Addr) {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	fmt.Println(newAddrs)
	for _, addr := range newAddrs {
		hosts.RemoveAddress(addr.host)
		hosts.AddHost(addr.ip.String(), addr.host)
	}
	hfData := hosts.RenderHostsFile()
	fmt.Println(hfData)
	hosts.Save()
	fmt.Println("save over")
}
