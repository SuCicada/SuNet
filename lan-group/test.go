package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main2() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	http.HandleFunc("/list",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r)
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})
	err := http.ListenAndServe(":414",nil)
	if err != nil {
		log.Panicln(err.Error())
	}
	fmt.Println("ok")
}
