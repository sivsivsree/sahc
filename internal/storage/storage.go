package storage

import (
	"encoding/json"
	"fmt"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
)

// Init is currently not using
func Init() {

}

// SaveConfigurations is used to save the conf state to leveldb
func SaveConfigurations(conf data.Configuration) error {
	db, err := leveldb.OpenFile(data.DB_NAME, nil)
	defer func(db *leveldb.DB) {
		err := db.Close()
		fmt.Println("SaveConfigurations, We are done", err)
	}(db)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("[SaveConfigurations]", conf)
	confByte, err := json.Marshal(conf);
	if err != nil {
		return err
	}

	return db.Put([]byte(data.CONF_KEY), confByte, nil)

}

// GetConfiguration is used to get the conf state from leveldb
func GetConfiguration() (*data.Configuration, error) {
	db, err := leveldb.OpenFile(data.DB_NAME, &opt.Options{ReadOnly: true})

	defer func(db *leveldb.DB) {
		err := db.Close()
		fmt.Println("GetConfiguration, We are done", err)
	}(db)

	if err != nil {
		return nil, err
	}

	confByte, err := db.Get([]byte(data.CONF_KEY), nil)
	if err != nil {
		return nil, err
	}
	var conf data.Configuration

	_ = json.Unmarshal(confByte, &conf)

	return &conf, nil
}
