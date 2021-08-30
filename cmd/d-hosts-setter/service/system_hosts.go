package service

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func SetSystemHosts(ipMapHostnames map[string]string, hostsPath string) error {

	hosts, err := readSystemHosts(hostsPath)
	if err != nil {
		return err
	}

	hostsStr := replaceIps(ipMapHostnames, hosts)

	err = writeSystemHosts(hostsPath, hostsStr)
	if err != nil {
		return err
	}

	return nil
}

func replaceIps(ipMapHostnames map[string]string, hosts string) string {
	hostsStr := hosts

	hostsArr := splitHosts(hostsStr)

	for ip, hostname := range ipMapHostnames {
		a := net.ParseIP(ip)
		if a != nil {
			ihs := ""
			has := false
			hasN := false

			for i, ipAndHostnames := range hostsArr {
				rIp, hasTemp := replaceIP(ip, hostname, ipAndHostnames)
				ihs = ihs + rIp
				if hasTemp {
					has = true
				}
				if i == len(hostsArr)-1 {
					if len(ipAndHostnames) > 0 {
						hasN = ipAndHostnames[len(ipAndHostnames)-1] == '\n'
					}
				}
			}

			if !has {
				if hasN {
					ihs = ihs + ip + " " + hostname + "\n"
				} else {
					ihs = ihs + "\n" + ip + " " + hostname
				}
			}

			hostsStr = ihs
			hostsArr = splitHosts(hostsStr)
		}
	}

	return hostsStr
}

func splitHosts(hosts string) []string {
	hostsArr := make([]string, 0)

	hostsArrTemp := strings.Split(hosts, "\n")
	for i, h := range hostsArrTemp {
		if i == len(hostsArrTemp)-1 {
			hostsArr = append(hostsArr, h)
		} else {
			hostsArr = append(hostsArr, h+"\n")
		}
	}
	return hostsArr
}

func replaceIP(ip, hostname, ipAndHostnames string) (string, bool) {
	ih := strings.TrimLeft(ipAndHostnames, " ")
	if len(ih) > 0 && ih[0] != '#' {
		if strings.Contains(ipAndHostnames, hostname) {
			spaces := getSpace(ipAndHostnames)
			ipEnd := strings.Index(ih, " ")
			if ipEnd > 0 {
				return spaces + ip + ih[ipEnd:], true
			}
		}
	}
	return ipAndHostnames, false
}

func getSpace(ih string) string {
	spaces := ""
	for _, s := range ih {
		if s != ' ' {
			return spaces
		}
		spaces += " "
	}
	return spaces
}

func writeSystemHosts(hostsPath, content string) error {
	stringBytes := []byte(content)
	if err := ioutil.WriteFile(hostsPath, stringBytes, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func readSystemHosts(hostsPath string) (string, error) {
	hosts := ""

	f, err := os.OpenFile(hostsPath, os.O_RDONLY, 0600)

	defer f.Close()

	if err != nil {
		return hosts, err
	} else {
		contentByte, err := ioutil.ReadAll(f)
		if err != nil {
			return hosts, err
		}
		return string(contentByte), nil
	}
}

func GetSystemDir() string {
	hostsPath := "/etc/hosts"
	if runtime.GOOS == "windows" {
		hostsPath = getWinSystemDir()
		hostsPath = filepath.Join(hostsPath, "system32", "drivers", "etc", "hosts")
	}
	return hostsPath
}

func getWinSystemDir() string {
	dir := ""
	if runtime.GOOS == "windows" {
		dir = os.Getenv("windir")
	}

	return dir
}
