package main

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var hostNameMapIp = make(map[string]string)

func getEtcdClient() (*clientv3.Client, error) {
	var etcdUrl = ""
	if os.Getenv("ETCDURL") != "" {
		etcdUrl = os.Getenv("ETCDURL")
	}

	log.Println("ETCDURL=" + etcdUrl)

	etcdUrlArr := strings.Split(etcdUrl, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdUrlArr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

type DnsValue struct {
	Host string `json:"host"`
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

	cli, err := getEtcdClient()
	defer cli.Close()
	if err != nil {
		log.Println(err)
	} else {
		hs := strings.Split(hostname, ".")
		reverseArray(hs)
		key := "/skydns/" + strings.Join(hs, "/") + "/"
		log.Println("key=" + key)
		log.Println("value=" + ip)
		valueByte, err := json.Marshal(DnsValue{
			Host: ip,
		})
		if err != nil {
			log.Println(err)
		} else {
			ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
			_, err := cli.Put(ctx, key, string(valueByte))
			if err != nil {
				log.Println(err)
			}
		}
	}

	_, _ = w.Write([]byte("hostname=" + hostname + ",ip=" + ip + "\n"))
}

func hHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`PUT /set
GET /get
`))
}

func reverseArray(numbers []string) {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
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
