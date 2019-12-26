package configurations

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func loadConfiguration(filename string) (*data.Configuration, error) {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	conf := &data.Configuration{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func getConfigPath() string {
	confPath := os.Getenv(data.ENV_CONFIG)
	return confPath
}

// initialize is the configuration
func initialize(db *leveldb.DB) error {

	conf, err := loadConfiguration(getConfigPath())

	if err != nil {
		log.Println("[Configuration file]", "is the \"SAHC_CONFIG\" env variable correct, currently pointing to config file '"+getConfigPath()+"'")
		return err
	}

	if err = conf.SaveConfigurations(db); err != nil {
		return err
	}

	return nil
}

// HotReload is used to periodically check the file changes.
func HotReload(change chan bool, db *leveldb.DB) chan bool {

	if err := initialize(db); err != nil {
		log.Fatal("[Configuration error]", err)
	}

	// How often to fire the passed in function
	// in milliseconds
	interval := 5 * time.Second

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)
	hash := make(chan string)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		prevHash, _ := md5sum(getConfigPath())
		for {

			pHash := prevHash
			select {
			case <-ticker.C:

				go func(pHash string, hash chan string) {

					newFileHash, _ := md5sum(getConfigPath())

					if newFileHash != pHash {
						hash <- newFileHash
					}

				}(pHash, hash)

			case <-clear:
				ticker.Stop()
				log.Println("[HotReload]", "Stopped")
				return

			case val := <-hash:

				log.Println("File change detected, updating configurations", val)
				if err := initialize(db); err != nil {
					log.Println("[Configuration error] ", err)
					clear <- true
				}
				prevHash = val
				// change should be triggered last.
				change <- true

			}

		}
	}()

	return clear

}

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}
