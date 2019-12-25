package sahc

import (
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/sivsivsree/sahc/internal/health"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"sync"
)

// Start is used to bootstrap Hot-Reload and Health together
func Start(db *leveldb.DB, restart chan bool) {
	// run the service runner to check if the services are running or not.
	activeSrvCh := make(chan data.HealthJobs, 100)
	var m sync.Mutex

	//var serviceRegistry []data.HealthJobs
	serviceRegistry := make(map[string]data.HealthJobs)

	go func(activeSrvCh chan data.HealthJobs, db *leveldb.DB) {
		var wg sync.WaitGroup
		for {
			select {
			case _ = <-restart:

				for key, rsv := range serviceRegistry {

					wg.Add(1)
					//go func(rsv data.HealthJobs, wg sync.WaitGroup, m sync.Mutex) {

					m.Lock()
					log.Println("[PAUSE] background scheduler", len(serviceRegistry))
					rsv.Running <- true
					delete(serviceRegistry, key)

					m.Unlock()
					wg.Done()
					//}(rsv, wg, m)

					wg.Wait()
				}

				log.Println("[START] Monitoring Health Started..")
				health.StartMonit(activeSrvCh, &m, db)

			case run := <-activeSrvCh:
				m.Lock()
				serviceRegistry[data.JobId(8)] = run
				m.Unlock()

			}
		}
	}(activeSrvCh, db)
}
