package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	ls := func(dir string) {
		files, _ := ioutil.ReadDir(dir)
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}

	file := os.Args[0]
	fmt.Println("current file: ", file)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println("current file in dir: ", dir)
	fmt.Println("now list files in current dir ")
	ls(dir)

	defer func() {
		err := os.Remove(file)
		if err != nil {
			println(err.Error())
		}
		fmt.Println("ok i delete ", file)
		fmt.Println("now let us show the dir")
		ls(dir)
	}()
}
