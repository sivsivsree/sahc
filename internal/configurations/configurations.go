package configurations

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/sivsivsree/sahc/internal/data"
	"github.com/sivsivsree/sahc/internal/storage"
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
	confPath := os.Getenv("SAHC_CONFIG")
	return confPath
}

// Initialize the configuration
func Init() error {

	conf, err := loadConfiguration(getConfigPath())

	if err != nil {
		return err
	}

	if err = storage.SaveConfigrations(*conf); err != nil {
		return err
	}

	return nil
}

// GetConfiguration is used to get the configuration from storage
func GetConfiguration() (*data.Configuration, error) {

	conf, err := storage.GetConfiguration()
	if err != nil {
		return conf, err
	}
	return nil, nil

}

func HotReload() chan bool {

	if err := Init(); err != nil {
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
				return

			case val := <-hash:
				log.Println("File change detected, updating configurations", val)
				if err := Init(); err != nil {
					fmt.Println("[Configuration error]", err)
					clear <- true
				}
				prevHash = val
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
