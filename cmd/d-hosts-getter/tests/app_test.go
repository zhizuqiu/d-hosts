package tests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"d-hosts/cmd/d-hosts-getter/model"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"
)

func TestEtcdCli(t *testing.T) {

	endpoints := []string{"etcd:2379"}

	var etcdCert = "certs/admin-zhizuqiu.pem"
	var etcdCertKey = "certs/admin-zhizuqiu-key.pem"
	var etcdCa = "certs/ca.pem"

	// 加载客户端证书
	cert, err := tls.LoadX509KeyPair(etcdCert, etcdCertKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 加载 CA 证书
	caData, err := ioutil.ReadFile(etcdCa)
	if err != nil {
		log.Fatal(err)
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	cfg := clientv3.Config{
		Endpoints: endpoints,
		TLS:       _tlsConfig, // Client.Config设置 TLS
	}

	cli, err := clientv3.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	hs := []string{"test"}
	ip := "127.0.0.2"
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
			return
		}
	}

	//  curl -X PUT "http://localhost:3000/set?hostname=test"
	resp, err := cli.Get(context.TODO(), "/skydns/test/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	// 关闭客户端连接
	defer cli.Close()
}
