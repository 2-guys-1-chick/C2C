package main

import (
	"log"

	"github.com/2-guys-1-chick/c2c/cfg"
	"github.com/2-guys-1-chick/c2c/network/client"
	"github.com/2-guys-1-chick/c2c/network/server"
)

func main() {
	err := server.StartServer(cfg.GetPort())
	if err != nil {
		log.Println(err)
	}

	connMgr := client.ConnManager{}
	connMgr.InitRoundup()

	quit := make(chan struct{})
	<-quit
}
