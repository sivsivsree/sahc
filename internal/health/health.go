package health

import (
	"fmt"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/sivsivsree/sahc/internal/storage"
	"log"
	"sync"
	"time"
)

// StartMonit will start monitoring the services defined

func StartMonit(runningSrv chan<- data.HealthJobs, m *sync.Mutex) {

	wg := sync.WaitGroup{}

	log.Println("[StartMonit]", "Load config")
	m.Lock()
	config, err := storage.GetConfiguration()

	if err != nil {
		log.Println("[StartMonit]", err)
	}

	for sid, service := range config.Services {

		wg.Add(1)

		go func(sid int, service data.Services) {
			running := runner(config, sid)
			run := data.HealthJobs{Running: running}
			runningSrv <- run
			//fmt.Println("runningSrv <- run")
			wg.Done()
		}(sid, service)

	}
	m.Unlock()
	wg.Wait()

}

func runner(config *data.Configuration, sid int) chan bool {

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
					fmt.Println("Status check | ", "Service ID:", sid, " |  Status:", service.Status, " |  Name:", service.Name)
					wg.Done()
				}(sid, service)

				wg.Wait()

			case <-clear:
				ticker.Stop()
				log.Println("[runner]", "Stopped", sid, service.Name, service.Status)
				close(clear)
				return

			}

		}
	}(clear)

	return clear

}

func checkStatus() bool {
	return true
}

// RemoveJob is used to remove sheduled health check jobs
func RemoveJob(svs []data.HealthJobs, index int) []data.HealthJobs {

	newJobList := []data.HealthJobs{}

	for i, s := range svs {
		if i != index {
			newJobList = append(newJobList, s)
		}
	}
	return newJobList
}

//func StopMonit() {
//
//}
