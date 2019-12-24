package cmd

import (
	"github.com/sivsivsree/sahc/internal/configurations"
	"github.com/sivsivsree/sahc/internal/health"
	"log"
	"os"
)

func main() {

	done := make(chan bool)
	// get configuration by flag or by ENV
	if configPath := os.Getenv("SAHC_CONFIG"); configPath == "" {
		log.Fatal("[configPath]", "No configurations passed")
	}


	// service to check the configuration changed or not.
	_, _ = configurations.GetConfiguration()
	configurations.HotReload()

	// run the service runner to check if the services are running or not.
	health.StartMonit()

	<-done
	// expose api to provide api interactions.

}
