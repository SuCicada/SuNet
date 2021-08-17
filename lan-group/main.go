package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var port = 4140

type Server interface {
}

func getSelfHost() string {
	return "win.ip"
}

var PLZ = []byte("PLZ")

func HttpServer(port int) {
	http.HandleFunc("/list",
		func(w http.ResponseWriter, r *http.Request) {
			s, _ := ioutil.ReadAll(r.Body)
			if bytes.Compare(s, PLZ) == 0 {
				log.Println("get", r.Header["User-Agent"], r.Method, r.RemoteAddr)
				res := fmt.Sprintf("OK:%s", getSelfHost())
				if _, err := fmt.Fprint(w, res); err != nil {
					log.Println(err.Error())
				}
			}
		})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panicln(err.Error())
	}
	fmt.Println("ok")
}


func main() {
	HttpServer(port)
	//HttpClient()
}
