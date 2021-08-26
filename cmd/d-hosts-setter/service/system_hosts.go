package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func SetSystemHosts(ip, hostname, hostsPath string) error {
	hosts, err := readSystemHosts(hostsPath)
	if err != nil {
		return err
	}

	ihs := ""

	hostsArr := strings.Split(hosts, "\n")
	for _, ipAndHostnames := range hostsArr {
		ihs = ihs + replaceIP(ip, hostname, ipAndHostnames) + "\n"
	}

	err = writeSystemHosts(hostsPath, ihs)
	if err != nil {
		return err
	}

	return nil
}

func replaceIP(ip, hostname, ipAndHostnames string) string {
	ih := strings.TrimLeft(ipAndHostnames, " ")
	if len(ih) > 0 && ih[0] != '#' {
		if strings.Contains(ipAndHostnames, hostname) {
			spaces := getSpace(ipAndHostnames)
			ipEnd := strings.Index(ih, " ")
			if ipEnd > 0 {
				return spaces + ip + ih[ipEnd:]
			}
		}
	}
	return ipAndHostnames
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
