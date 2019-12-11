// +build static
// +build !nostatic !cli

package main

import (
	"log"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"time"
)

import (
	. "github.com/eyedeekay/sam-forwarder/gui"
	"github.com/zserge/lorca"
)

var view lorca.UI

func RunUI() {
	view, err = LaunchUI(&s)
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *App) Serve() bool {
	log.Println("Starting Tunnels()")
	for _, element := range s.clientMux.Tunnels() {
		log.Println("Starting service tunnel", element.ID())
		go element.Serve()
	}

	if s.UseWebUI() == true {
		go s.clientMux.ListenAndServe()
		if view, err = LaunchUI(s); err != nil {
			log.Println(err.Error())
			return Exit()
		} else {
			return Exit()
		}
	} else {
		return Exit()
	}
	return false
}

func Exit() bool {
	Close := false
	for !Close {
		time.Sleep(1 * time.Second)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				log.Println(sig)
				if view != nil {
					view.Close()
				}
				Close = true
			}
		}()
	}
	return false
}

func FixDir() {
	os.MkdirAll(*userName, 0755)
	if u, err := user.Lookup(*chromeUser); err != nil {
		switch err.(type) {
		case user.UnknownUserError:
			log.Println(u, "user not found")
		}
	} else if u, err := user.Lookup(*userName); err != nil {
		switch err.(type) {
		case user.UnknownUserError:
			log.Println(u, "user not found")
		}
	} else {
		ug, _ := strconv.Atoi(u.Uid)
		gg, _ := strconv.Atoi(u.Gid)

		ChownR(*userName, ug, gg)
	}
}
