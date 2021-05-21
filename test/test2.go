package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func GetCmdNews(v1 string) string {
	//cmd := exec.Command(v1)
	cmd := exec.Command("/bin/bash", "-c", v1)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		//	log.Fatal(err)
		return ""
	}

	if err := cmd.Wait(); err != nil {
		//log.Fatal(err)
		return ""
	}
	getVersion := fmt.Sprintf("%s", string(opBytes))
	return getVersion
}

func main() {

	confFile := `find ../../zyf-jeecg-boot/ -type f -iname "*config*" -o -iname "web.xml" -o -iname "*database*" -o -iname "*pass*" 2>/dev/null`
	result := GetCmdNews(confFile)
	fmt.Println(result)
}
