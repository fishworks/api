package main

import (
	"log"

	"github.com/fishworks/api/server"
)

const (
	// Usage is the usage string of this client, which is displayed on os.Stdout when requested.
	Usage string = "usage: api <proto://host:port>"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("deis[api]: ")

	server, err := server.NewServer("unix", "/var/run/api.sock")
	if err != nil {
		log.Fatal(err)
	}
	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
