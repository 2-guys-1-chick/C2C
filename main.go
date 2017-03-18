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

	quit := make(chan struct{})
	collector := datcol.Collector{}
	collector.SetDistributor(srv)
	go collector.Run(quit)

	repHandler := datrep.InitHandler()

	connMgr := client.ConnManager{}
	connMgr.SetPacketHandler(repHandler)

	connMgr.Connect("localhost", cfg.GetPort())
	connMgr.InitRoundup()

	<-quit
}
