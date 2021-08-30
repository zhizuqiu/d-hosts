package service

import (
	"fmt"
	"testing"
)

var hostsPath = ""

func init() {
	hostsPath = GetSystemDir()
}

func TestReadSystemHosts(t *testing.T) {
	s, _ := readSystemHosts(hostsPath)
	fmt.Println(s)
}

func TestReplaceIP(t *testing.T) {
	ip := "127.0.0.1"
	hostname := "hostname_test"

	var replaceIPTests = []struct {
		in       string
		expected string
	}{
		{"   127.0.0.2    hostname_test     ", "   127.0.0.1    hostname_test     "},
		{"127.0.0.2 hostname_test", "127.0.0.1 hostname_test"},
		{"#127.0.0.2 hostname_test", "#127.0.0.2 hostname_test"},
		{"   #127.0.0.2 hostname_test", "   #127.0.0.2 hostname_test"},
		{"", ""},
		{"127.0.0.2    hostname_test", "127.0.0.1    hostname_test"},
		{"127.0.0.2    hostname_test     ", "127.0.0.1    hostname_test     "},
		{"127.0.0.2 hostname_test\n", "127.0.0.1 hostname_test\n"},
	}

	for _, tt := range replaceIPTests {
		actual, _ := replaceIP(ip, hostname, tt.in)
		if actual != tt.expected {
			t.Errorf("replaceIP(%s) = \"%s\"; expected \"%s\"", tt.in, actual, tt.expected)
		}
	}
}

func TestReplaceIps(t *testing.T) {

	ipMapHostnames := make(map[string]string)
	ipMapHostnames["172.0.0.1"] = "host1"
	ipMapHostnames["172.0.0.2"] = "host2"

	var replaceIPsTests = []struct {
		in       string
		expected string
	}{
		{"", `
172.0.0.1 host1
172.0.0.2 host2`},
		{"172.0.0.1 host1", `172.0.0.1 host1
172.0.0.2 host2`},
		{`172.0.0.1 host1
172.0.0.2 host2`, `172.0.0.1 host1
172.0.0.2 host2`},
		{`172.0.0.3 host1
172.0.0.4 host2`, `172.0.0.1 host1
172.0.0.2 host2`},
	}

	for _, tt := range replaceIPsTests {
		actual := replaceIps(ipMapHostnames, tt.in)
		if actual != tt.expected {
			t.Errorf("\nreplaceIP(%s) = \n\"%s\"\nexpected: \n\"%s\"", tt.in, actual, tt.expected)
		}
	}

}
