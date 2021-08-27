package main

import (
	"log"
	"net/http"
	"strings"
)

var hostNameMapIp = make(map[string]string)

func hHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`PUT /set
GET /get
`))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("405 Method Not Allowed"))
		return
	}

	ip := parseIp(r.RemoteAddr)
	ipQuery := r.URL.Query().Get("ip")
	if ipQuery != "" {
		ip = ipQuery
	}

	hostname := r.URL.Query().Get("hostname")

	hostNameMapIp[hostname] = ip

	_, _ = w.Write([]byte("hostname=" + hostname + ",ip=" + ip + "\n"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	h := ""
	for hostname, ip := range hostNameMapIp {
		h += hostname + " " + ip + "\n"
	}
	_, _ = w.Write([]byte(h))
}

func parseIp(remoteAddr string) string {
	ip := remoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index > 0 {
		ip = remoteAddr[:index]
	}
	return ip
}

func main() {
	h := http.HandlerFunc(hHandler)
	set := http.HandlerFunc(setHandler)
	get := http.HandlerFunc(getHandler)

	http.Handle("/", h)
	http.Handle("/set", set)
	http.Handle("/get", get)

	log.Println("Listening 3000... ")
	_ = http.ListenAndServe(":3000", nil)
}
