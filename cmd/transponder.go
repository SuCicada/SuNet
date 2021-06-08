package main

import (
	"flag"
	"github.com/SuCicada/SuNet/transponder"
	"strconv"
	"strings"
)

func main() {
	mode := flag.String("mode", "", "select mode (client/server)")
	remoteServerAddr := flag.String("f", "", "remote host:listenPort")
	localPort := flag.Int("lp", -1, "lcoalPort")
	remoteTransPort := flag.Int("rp", -1, "remoteTransPort")

	listenPort := flag.Int("p", -1, "listenPort")

	flag.Parse()
	switch *mode {
	case "c", "client":
		if *remoteServerAddr == "" {
			flag.Usage()
			return
		}
		remoteIp, remotePort := func() (string, int) {
			sp := strings.Split(*remoteServerAddr, ":")
			n, err := strconv.Atoi(sp[1])
			if err != nil {
				flag.Usage()
			}
			return sp[0], n
		}()
		transponder.Client(remoteIp, remotePort, *localPort, *remoteTransPort)
	case "s", "server":
		if *listenPort == -1 {
			flag.Usage()
			return
		}
		transponder.Server(*listenPort)
	default:
		flag.Usage()
	}
}
