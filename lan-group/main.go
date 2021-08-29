package lan_group

import (
	"flag"
)

var port int
var selfHost string

func getSelfHost() string {
	return selfHost
}

var PLZ = []byte("PLZ")

func Main() {
	writeFlag := flag.Bool("w", false, "write hosts file ?")
	serverFlag := flag.Bool("s", false, "open the listening server ?")
	flag.IntVar(&port, "p", 4140, "port")
	flag.StringVar(&selfHost, "h", "", "self host name ")
	flag.Parse()
	SelfInfo()
	if *serverFlag {
		if selfHost == "" {
			flag.Usage()
			return
		}
		HttpServer(port)
	} else {
		HttpClient(*writeFlag)
	}
}
