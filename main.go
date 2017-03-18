package main

import (
	"log"

	"github.com/2-guys-1-chick/c2c/cfg"
	"github.com/2-guys-1-chick/c2c/datcol"
	"github.com/2-guys-1-chick/c2c/datrep"
	"github.com/2-guys-1-chick/c2c/network/client"
	"github.com/2-guys-1-chick/c2c/network/server"
)

func main() {
	srv, err := server.StartServer(cfg.GetPort())
	if err != nil {
		log.Println(err)
	}

	collector := datcol.Collector{}
	collector.SetDistributor(srv)
	go collector.Run()

	repHandler := datrep.InitHandler()

	connMgr := client.ConnManager{}
	connMgr.SetPacketHandler(repHandler)
	connMgr.InitRoundup()

	quit := make(chan struct{})
	<-quit
}
