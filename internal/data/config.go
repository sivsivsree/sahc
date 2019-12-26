package data

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

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

// UpdateStatus is used to update the health availability status of service
func (conf *Configuration) UpdateStatus(db *leveldb.DB, id int, status bool) error {

	conf.Services[id].Status = status

	// log.Println("[UpdateStatus]", conf)
	confByte, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return db.Put([]byte(CONF_KEY), confByte, nil)
}

// SaveConfigurations is used to Save the conf instance to leveldb
func (conf *Configuration) SaveConfigurations(db *leveldb.DB) error {

	log.Println("[SaveConfigurations]", "Saving..")
	confByte, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return db.Put([]byte(CONF_KEY), confByte, nil)

}

// GetConfiguration is used to get configuration from leveldb
func (conf Configuration) GetConfiguration(db *leveldb.DB) (*Configuration, error) {

	//defer func(db *leveldb.DB) {
	//	//err := db.Close()
	//
	//}(db)

	confByte, err := db.Get([]byte(CONF_KEY), nil)
	if err != nil {
		return nil, err
	}
	var c Configuration

	_ = json.Unmarshal(confByte, &c)

	return &c, nil
}
