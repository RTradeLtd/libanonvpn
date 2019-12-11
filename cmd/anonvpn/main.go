package main

import (
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

import (
	samtunnelhandler "github.com/RTradeLtd/go-anonvpn/ctrl"
	i2pbrowserproxy "github.com/eyedeekay/httptunnel/multiproxy"
)

type App struct {
	clientMux *samtunnelhandler.TunnelHandlerMux
	eepProxy  *i2pbrowserproxy.SAMMultiProxy
	Close     bool
}

func (s *App) UseWebUI() bool {
	return *webAdmin
}

func (s *App) Title() string {
	return *userName
}

func (s *App) Width() int {
	return 800
}

func (s *App) Height() int {
	return 600
}

func (s *App) Resizable() bool {
	return true
}

func (s *App) URL() string {
	return "http://" + webHost + ":" + *webPort
}

var webHost = "127.0.0.1"

func User() string {
	runningUser, _ := user.Current()
	if runtime.GOOS != "windows" {
		if runningUser.Uid == "0" {
			cmd := exec.Command("logname")
			out, err := cmd.Output()
			if err != nil {
				return err.Error()
			}
			return string(out)
		}
	}
	return runningUser.Name
}

var runningUser = User()

var (
	s          App
	err        error
	accessList flagOpts
)

type flagOpts []string

func (f *flagOpts) String() string {
	r := ""
	for _, s := range *f {
		r += s + ","
	}
	return strings.TrimSuffix(r, ",")
}

func (f *flagOpts) Set(s string) error {
	*f = append(*f, s)
	return nil
}

func (f *flagOpts) StringSlice() []string {
	var r []string
	for _, s := range *f {
		r = append(r, s)
	}
	return r
}
