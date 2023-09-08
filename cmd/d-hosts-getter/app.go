package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"d-hosts/cmd/d-hosts-getter/model"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var hostNameMapIp = make(map[string]string)

func getEtcdClient() (*clientv3.Client, error) {
	var etcdUrl = "https://etcd:2379"
	if os.Getenv("ETCDURL") != "" {
		etcdUrl = os.Getenv("ETCDURL")
	}
	log.Println("ETCDURL=" + etcdUrl)

	var etcdCert = "certs/admin-zhizuqiu.pem"
	if os.Getenv("ETCDCERT") != "" {
		etcdCert = os.Getenv("ETCDCERT")
	}
	log.Println("ETCDCERT=" + etcdCert)
	var etcdCertKey = "certs/admin-zhizuqiu-key.pem"
	if os.Getenv("ETCDCERTKEY") != "" {
		etcdCertKey = os.Getenv("ETCDCERTKEY")
	}
	log.Println("ETCDCERTKEY=" + etcdCertKey)
	var etcdCa = "certs/ca.pem"
	if os.Getenv("ETCDCA") != "" {
		etcdCa = os.Getenv("ETCDCA")
	}
	log.Println("ETCDCA=" + etcdCa)

	// 加载客户端证书
	cert, err := tls.LoadX509KeyPair(etcdCert, etcdCertKey)
	if err != nil {
		return nil, err
	}

	// 加载 CA 证书
	caData, err := ioutil.ReadFile(etcdCa)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	etcdUrlArr := strings.Split(etcdUrl, ",")

	cfg := clientv3.Config{
		Endpoints:   etcdUrlArr,
		DialTimeout: 5 * time.Second,
		TLS:         _tlsConfig,
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return cli, nil
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
		valueByte, err := json.Marshal(model.DnsValue{
			Host: ip,
		})
		log.Println("value=" + string(valueByte))
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

	var listenAddr = ":31006"
	if os.Getenv("Listen_ADDR") != "" {
		listenAddr = os.Getenv("Listen_ADDR")
	}
	log.Println("Listen_ADDR=" + listenAddr)

	http.Handle("/", h)
	http.Handle("/set", set)
	http.Handle("/get", get)

	log.Println("Listening " + listenAddr + "... ")
	_ = http.ListenAndServe(listenAddr, nil)
}
