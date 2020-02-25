// +build cli
// +build !nostatic !static

package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func RunUI() {

}

func (s *App) Serve() bool {
	log.Println("Starting Tunnels()")
	for _, element := range s.clientMux.Tunnels() {
		log.Println("Starting service tunnel", element.ID())
		go element.Serve()
	}

	if err := Canal(); err != nil {
		return false
	}

	return Exit()
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

				Close = true
			}
		}()
	}
	return false
}

func FixDir() {

}
