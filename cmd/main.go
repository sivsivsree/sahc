package main

import (
	"context"
	"github.com/sivsivsree/sahc/internal/configurations"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/sivsivsree/sahc/internal/sahc"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	db, err := leveldb.OpenFile(data.DB_NAME, nil)

	if err != nil {
		log.Fatal(err)

	}
	// for graceful shutdown of service.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// restart on conf reload
	restart := make(chan bool, 30)

	// service to check the configuration changed or not.
	clear := configurations.HotReload(restart, db)

	sahc.Start(db, restart)

	// need to kickstart the service initally
	restart <- true

	<-done

	log.Println("Releasing all the allocated resources")

	// Stop the HotReload

	// Gracefull Shutdown added.
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		clear <- true
		_ = db.Close()
		// extra handling here
		cancel()
	}()

}
