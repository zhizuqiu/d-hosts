package service

import (
	"fmt"
	"testing"
)

var hostsPath = ""

func init() {
	hostsPath = GetSystemDir()
}

func TestSetSystemHosts(t *testing.T) {
	_ = SetSystemHosts("127.0.0.1", "hostname_test", hostsPath)
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
	}

	for _, tt := range replaceIPTests {
		actual := replaceIP(ip, hostname, tt.in)
		if actual != tt.expected {
			t.Errorf("replaceIP(%s) = \"%s\"; expected \"%s\"", tt.in, actual, tt.expected)
		}
	}
}
