package service

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}}

}

func GetIp(address string) (string, error) {
	resp, err := client.Get(address + "/ip")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
