package main

import (
	"flag"
	"fmt"
	"github.com/SuCicada/SuNet/process"
	"log"
)

func main() {
	port := flag.Int("p", -1, "port")
	isKill := flag.Bool("k", false, "is kill")
	flag.Parse()
	if *port <= 0 {
		flag.Usage()
		return
	}
	pid := process.FindPortPid(*port)
	fmt.Println("pid", pid)
	if *isKill {
		process.KillPid(pid)
		log.Println("kill", pid)
	}
}
