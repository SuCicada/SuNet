package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
)

func readFile(file string) string {
	data, _ := ioutil.ReadFile(file)
	return string(data)
}

func test() {
	user, _ := user.Current()
	files := []string{
		user.HomeDir + "/.bash_history",
		"/etc/resolv.conf",
		"/etc/*release",
	}
	for _, file := range files {
		res := readFile(file)
		println(file)
		println(res)
	}
	//println(os.Getenv("$HISTFILE"))
	//println(os.Getenv("HISTFILE"))
}

// 这里为了简化，我省去了stderr和其他信息
func Command(cmd string) error {
	println(cmd)
	c := exec.Command("bash", "-c", cmd)
	// 此处是windows版本
	// c := exec.Command("cmd", "/C", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}
