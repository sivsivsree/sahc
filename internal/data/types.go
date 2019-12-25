package data

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
)

type HealthJobs struct {
	Running chan bool
}

// Services is used to store service
type Services struct {
	Name     string `yaml:"name";json:"name"`
	Interval int    `yaml:"interval";json:"interval"`
	Status   bool   `json:"status"`
}

// Configuration is used to parse and use config.yaml
type Configuration struct {
	Version  float64    `yaml:"version";json:"version"`
	Services []Services `yaml:"services";json:"services"`
}

func (conf *Configuration) UpdateStatus(db *leveldb.DB, id int, status bool) error {

	conf.Services[id].Status = status

	// log.Println("[UpdateStatus]", conf)
	confByte, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return db.Put([]byte(CONF_KEY), confByte, nil)
}
