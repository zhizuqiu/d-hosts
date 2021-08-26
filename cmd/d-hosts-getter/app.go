package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var remoteAddr = ""

func hHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`POST /curl
GET /get
GET /ip
`))
}

func curlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("405 Method Not Allowed"))
		return
	}
	fmt.Println(r)
	remoteAddr = r.RemoteAddr
	_, _ = w.Write([]byte("The Router Addr is: " + remoteAddr))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("The Router Addr is: " + remoteAddr))
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	ip := remoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index > 0 {
		ip = remoteAddr[:index]
	}
	_, _ = w.Write([]byte(ip))
}

func main() {
	h := http.HandlerFunc(hHandler)
	ch := http.HandlerFunc(curlHandler)
	get := http.HandlerFunc(getHandler)
	ip := http.HandlerFunc(ipHandler)

	http.Handle("/", h)
	http.Handle("/curl", ch)
	http.Handle("/get", get)
	http.Handle("/ip", ip)

	log.Println("Listening 3000... ")
	_ = http.ListenAndServe(":3000", nil)
}
