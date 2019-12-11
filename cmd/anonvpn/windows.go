// +build windows

package main

import (
	//"context"
	"flag"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.Run()
	return nil
}

func (p *program) Run() {
	lbMain(nil)
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "go-anonvpn",
		DisplayName: "AnonVPN Service",
		Description: "A VPN Service that uses I2P",
	}

	flag.Var(&accessList, "accesslistmembers", "Specify an access list member(can be used multiple times).")
	flag.Parse()

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		panic(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
