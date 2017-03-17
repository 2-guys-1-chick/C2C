package main

import (
	"github.com/2-guys-1-chick/c2c/cfg"
	"github.com/2-guys-1-chick/c2c/network/client"
	"github.com/2-guys-1-chick/c2c/network/server"
)

func main() {
	server.StartServer(cfg.GetPort())

	client.RoundupConnect(nil)

	quit := make(chan struct{})
	<-quit
}
