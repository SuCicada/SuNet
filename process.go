package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("p", -1, "port")
	flag.Parse()
	if *port <= 0 {
		flag.Usage()
		return
	}
	fmt.Println(*port)
	//pid := process.FindPortPid(port)
	//process.KillPid(pid)
}
