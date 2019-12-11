// +build !windows
// +build !android

package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"strings"
)

import (
	"crawshaw.io/littleboss"
	//	"github.com/eyedeekay/portcheck"
	"github.com/eyedeekay/sam-forwarder/hashhash"
)

func main() {
	lb := littleboss.New("name")
	lb.Run(func(ctx context.Context) {
		flag.Var(&accessList, "accesslistmembers", "Specify an access list member(can be used multiple times).")
		flag.Parse()

		FixDir()

		if *peoplehash != "" {
			slice := strings.Split(*peoplehash, " ")
			if length, err := strconv.Atoi(slice[len(slice)-1]); err == nil {
				Hasher, err := hashhash.NewHasher(length)
				if err != nil {
					return
				}
				lhash, err := Hasher.Unfriendlyslice(slice[0 : len(slice)-2])
				if err != nil {
					return
				}
				log.Println(lhash + ".b32.i2p")
			} else {
				Hasher, err := hashhash.NewHasher(52)
				if err != nil {
					return
				}
				lhash, err := Hasher.Unfriendlyslice(slice)
				if err != nil {
					return
				}
				log.Println(lhash + ".b32.i2p")
			}
			return
		}
		lbMain(ctx)
	})
}
