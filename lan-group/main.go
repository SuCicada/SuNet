package lan_group

import (
	"flag"
)

var selfHost string
var port int

func getSelfHost() string {
	return selfHost
}

var PLZ = []byte("PLZ")
var debugIP = ""
var scanScope = ""

func Main() {
	serverFlag := flag.Bool("s", false, "open the listening server ?")
	writeFlag := flag.Bool("w", false, "write hosts file ?")
	flag.IntVar(&port, "p", 4140, "port")
	flag.StringVar(&selfHost, "h", "", "self host name ")
	flag.StringVar(&debugIP, "debug", "", "debug one ip")
	flag.StringVar(&scanScope, "scan", "", "scan scope ip. eg: 10.0.3.00/24")
	flag.Parse()

	SelfInfo()

	if *serverFlag {
		if selfHost == "" {
			flag.Usage()
			return
		}
		HttpServer(port)
	} else {
		HttpClient(*writeFlag, scanScope)
	}
}
