package main

import (
	"github.com/sivsivsree/sahc/internal/configurations"
	"github.com/sivsivsree/sahc/internal/health"
	"log"
)

func main() {

	// get configuration by flag or by ENV
	if err := configurations.Init(); err != nil {
		log.Fatal("[configurations Init]", err)
	}

	// service to check the configuration changed or not.
	configurations.HotReload()

	// run the service runner to check if the services are running or not.
	health.StartMonit()

	// expose api to provide api interactions.

}
