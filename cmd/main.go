package main

import (
	"flag"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fishworks/api/server"
	"github.com/fishworks/api/settings"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/rest"
)

func validateSettings() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := kubernetes.NewForConfig(config); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&settings.ListenAddress, "a", "tcp://0.0.0.0:8080", "")
	flag.StringVar(&settings.ListenAddress, "addr", "tcp://0.0.0.0:8080", "")
	flag.StringVar(&settings.LogLevel, "l", "info", "")
	flag.StringVar(&settings.LogLevel, "log-level", "info", "")
	flag.Parse()

	if level, err := log.ParseLevel(settings.LogLevel); err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
	}
	validateSettings()

	protoAndAddr := strings.SplitN(settings.ListenAddress, "://", 2)
	server, err := server.New(protoAndAddr[0], protoAndAddr[1])
	if err != nil {
		log.Fatalf("failed to create server at %s: %v", settings.ListenAddress, err)
	}
	log.Printf("server is now listening at %s", settings.ListenAddress)
	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
