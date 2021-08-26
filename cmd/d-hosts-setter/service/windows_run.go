package service

import (
	"fmt"
	"github.com/kardianos/service"
	"log"
	"net"
	"os"
	"time"
)

var serviceConfig = &service.Config{
	Name:        "Hosts Setter",
	DisplayName: "Hosts Setter",
	Description: "通过访问远程d-hosts-getter接口，动态更新本地hosts文件",
}

func WindowsRun(address, hostname string, interval int) {

	prog := &Program{
		hostsPath: GetSystemDir(),
		address:   address,
		hostname:  hostname,
		interval:  interval,
	}
	s, err := service.New(prog, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[1]

	if cmd == "install" {
		err = s.Install()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("安装成功")
	} else if cmd == "uninstall" {
		err = s.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("卸载成功")
	} else {
		err = s.Run()
		if err != nil {
			_ = logger.Error(err)
		}
		return
	}

}

type Program struct {
	hostsPath string
	address   string
	hostname  string
	interval  int
}

func (p *Program) Start(s service.Service) error {
	log.Println("开始服务")
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	log.Println("停止服务")
	return nil
}

func (p *Program) run() {

	for {
		ip, err := GetIp(p.address)
		if err != nil {
			log.Println(err)
		} else {
			a := net.ParseIP(ip)
			if a == nil {
				log.Println("未获取到ip地址")
			} else {
				err = SetSystemHosts(a.String(), p.hostname, p.hostsPath)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("更新 hosts 成功")
				}
			}
		}

		time.Sleep(time.Second * time.Duration(p.interval))
	}
}
