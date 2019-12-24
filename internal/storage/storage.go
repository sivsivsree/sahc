package storage

import (
	"encoding/json"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

func Init() {

}

func SaveConfigrations(conf data.Configuration) error {
	db, err := leveldb.OpenFile(data.DB_NAME, nil)
	defer db.Close();
	if err != nil {
		return err
	}
	log.Println("[SaveConfigrations]", conf)
	confByte, err := json.Marshal(conf);
	if err != nil {
		return err
	}
	return db.Put([]byte(data.CONF_KEY), confByte, nil)

}

func GetConfiguration() (*data.Configuration, error) {
	db, err := leveldb.OpenFile(data.DB_NAME, nil)
	defer db.Close();
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
