package service

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
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

func GetIpMapHostnames(address string) (map[string]string, error) {

	resp, err := client.Get(address + "/get")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var ipMapHostnames IpMapHostnames
	ipMapHostnameList := strings.Split(string(body), "\n")
	for _, ipMapHostname := range ipMapHostnameList {
		ih := strings.Split(ipMapHostname, " ")
		if len(ih) > 1 {
			ipMapHostnames = ipMapHostnames.Append(ih[1], ih[0])
		}
	}

	return ipMapHostnames.ToIpMapString(), nil
}

type IpMapHostnames map[string][]string

func (i IpMapHostnames) Append(ip, hostname string) IpMapHostnames {
	if i == nil {
		i = make(map[string][]string)
	}

	h, ok := i[ip]
	if ok {
		h = append(h, hostname)
		i[ip] = h
	} else {
		i[ip] = []string{hostname}
	}

	return i
}

func (i IpMapHostnames) ToIpMapString() map[string]string {
	ipMapString := make(map[string]string)

	for ip, hostnames := range i {
		ipMapString[ip] = strings.Join(hostnames, " ")
	}

	return ipMapString
}
