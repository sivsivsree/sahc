package main

import (
	"context"
	"github.com/sivsivsree/sahc/internal/configurations"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/sivsivsree/sahc/internal/health"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"os/signal"
	"sync"
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

	// run the service runner to check if the services are running or not.
	runingSrvChan := make(chan data.HealthJobs, 100)
	var m sync.Mutex

	//var runingSrvs []data.HealthJobs
	runingSrvs := make(map[string]data.HealthJobs)

	go func(runingSrvChan chan data.HealthJobs, db *leveldb.DB) {

		var wg sync.WaitGroup
		for {
			select {
			case _ = <-restart:

				for key, rsv := range runingSrvs {

					wg.Add(1)
					//go func(rsv data.HealthJobs, wg sync.WaitGroup, m sync.Mutex) {

					m.Lock()
					log.Println("stop the service..", len(runingSrvs), runingSrvs)
					rsv.Running <- true
					delete(runingSrvs, key)
					log.Println("running service after remove..", len(runingSrvs), runingSrvs)
					m.Unlock()
					wg.Done()
					//}(rsv, wg, m)

					wg.Wait()
				}

				log.Println("Health Start Monitoring..")
				health.StartMonit(runingSrvChan, &m, db)

			case run := <-runingSrvChan:
				m.Lock()
				runingSrvs[data.JobId(8)] = run
				m.Unlock()

			}
		}
	}(runingSrvChan, db)

	restart <- true

	<-done

	log.Println("Releasing all the allocated resources")
	db.Close()
	// Stop the HotReload
	clear <- true
	// Gracefull Shutdown added.
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {

		// extra handling here
		cancel()
	}()

}
