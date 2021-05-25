package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	//confFile := `find ../../zyf-jeecg-boot/ -type f -iname "*config*" -o -iname "web.xml" -o -iname "*database*" -o -iname "*pass*" 2>/dev/null`
	//result := GetCmdNews(confFile)
	//fmt.Println(result)
	os.Args = strings.Split(". -r ubuntu.wsl2:22", " ")
	target := flag.String("r", "??", "the remote server (<host>:<port>)")
	b := func() *string {
		a := new(string)
		*a = "243242"
		return a
	}

	flag.Parse()
	fmt.Println(target)
}
