package main

import (
	"flag"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fishworks/api/server"
)

var (
	addrFlag string
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	flag.StringVar(&addrFlag, "a", "tcp://0.0.0.0:8080", "")
	flag.StringVar(&addrFlag, "addr", "tcp://0.0.0.0:8080", "")
	flag.Parse()

	protoAndAddr := strings.SplitN(addrFlag, "://", 2)
	server, err := server.NewServer(protoAndAddr[0], protoAndAddr[1])
	if err != nil {
		log.Fatalf("failed to create server at %s: %v", addrFlag, err)
	}
	log.Printf("server is now listening at %s", addrFlag)
	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
