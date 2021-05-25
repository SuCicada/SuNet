// +build windows

/*
	杀死 windows 下占用端口的可恶进程
*/

package process

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func runCmd(cmd string) string {
	exe := exec.Command("cmd.exe", "/C", cmd)
	res, err := exe.CombinedOutput()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return string(res)
}
func FindPortPid(port int) int {
	// netstat -aon|findstr "8080"
	res := runCmd(fmt.Sprintf("netstat -aon|findstr %d", port))

	var pid string
	for _, s := range strings.Split(strings.TrimSpace(res), "\n") {
		fields := strings.Fields(s)
		if strings.Contains(fields[1], strconv.Itoa(port)) {
			pid = fields[4]
			break
		}
	}
	p, _ := strconv.Atoi(pid)
	return p
}

func KillPid(pid int) {
	// taskkill /T /F /PID
	runCmd(fmt.Sprintf("taskkill /T /F /PID %d", pid))
}
