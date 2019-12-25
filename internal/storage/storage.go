package storage

import (
	"encoding/json"
	"fmt"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

// SaveConfigurations is used to save the conf state to leveldb
func SaveConfigurations(conf data.Configuration, db *leveldb.DB) error {

	log.Println("[SaveConfigurations]", conf)
	confByte, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return db.Put([]byte(data.CONF_KEY), confByte, nil)

}

// GetConfiguration is used to get the conf state from leveldb
func GetConfiguration(db *leveldb.DB) (*data.Configuration, error) {

	defer func(db *leveldb.DB) {
		//err := db.Close()
		fmt.Println("GetConfiguration, We are done")
	}(db)

	confByte, err := db.Get([]byte(data.CONF_KEY), nil)
	if err != nil {
		return nil, err
	}
	var conf data.Configuration

	_ = json.Unmarshal(confByte, &conf)

	return &conf, nil
}
