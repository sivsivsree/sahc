package health

import (
	"fmt"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net"
	"sync"
	"time"
)

// StartMonit will start monitoring the services defined
func StartMonit(runningSrv chan<- data.HealthJobs, m *sync.Mutex, db *leveldb.DB) {

	wg := sync.WaitGroup{}

	log.Println("[StartMonit]", "Load config")
	m.Lock()

	config, err := (&data.Configuration{}).GetConfiguration(db)

	if err != nil {
		log.Println("[StartMonit]", err)
	}

	for sid, service := range config.Services {

		wg.Add(1)

		go func(sid int, service data.Services, db *leveldb.DB) {
			running := runner(config, sid, db)
			run := data.HealthJobs{Running: running}
			runningSrv <- run
			//fmt.Println("runningSrv <- run")
			wg.Done()
		}(sid, service, db)

	}
	m.Unlock()
	wg.Wait()

}

func runner(config *data.Configuration, sid int, db *leveldb.DB) chan bool {

	service := config.Services[sid]

	interval := time.Duration(service.Interval) * time.Second

	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	go func(clear chan bool) {
		for {

			select {
			case <-ticker.C:

				// go checkSomething
				wg := sync.WaitGroup{}
				go func(sid int, services data.Services) {
					wg.Add(1)
					go statusCheckAndUpdate(config, sid, &wg, db)

				}(sid, service)

				wg.Wait()

			case <-clear:
				ticker.Stop()
				log.Println("[STOP]", "background scheduler", sid, service.Name, service.Status)
				close(clear)
				return

			}

		}
	}(clear)

	return clear

}

func statusCheckAndUpdate(conf *data.Configuration, sid int, wg *sync.WaitGroup, db *leveldb.DB) {

	timeout := time.Second
	conn, _ := net.DialTimeout("tcp", conf.Services[sid].Name, timeout)

	if conn != nil {
		defer conn.Close()
	}
	if err := conf.UpdateStatus(db, sid, conn != nil); err != nil {
		fmt.Println("Error", err)
	}
	log.Println("[Status check] ", "Service ID:", sid, " |  Status:", conf.Services[sid].Status, " |  Name:", conf.Services[sid].Name)

	wg.Done()
}
