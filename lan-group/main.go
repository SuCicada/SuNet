package main

import (
	"flag"
)

var port int
var selfHost string

func getSelfHost() string {
	return selfHost
}

var PLZ = []byte("PLZ")

func main() {
	writeFlag := flag.Bool("w", false, "write hosts file ?")
	serverFlag := flag.Bool("s", false, "open the listening server ?")
	flag.IntVar(&port, "s", 4140, "open the listening server ?")
	flag.StringVar(&selfHost, "h", "", "open the listening server ?")
	flag.Parse()
	SelfInfo()
	if *serverFlag {
		HttpServer(port)
	} else {
		HttpClient(*writeFlag)
	}
}
