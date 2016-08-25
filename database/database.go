package database

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/superdecimal/dolce/config"
	"github.com/superdecimal/dolce/logbook"
)

type Database struct {
	DatabaseName string
	Filename     string
	Path         string
	File         *os.File
	Version      int
	Data         map[string][]byte
	dbMutex      sync.Mutex
	Dlog         logbook.Logbook
}

func init() {
}

// Set adds a value to the map and stores the action on the log
func (d *Database) Set(key string, value string) error {
	d.dbMutex.Lock()
	defer d.dbMutex.Unlock()

	data := []byte(value)
	d.Dlog.Append(key, data)
	d.Data[key] = data

	return nil
}

func (d *Database) Read(key string) (string, error) {
	return string(d.Data[key]), nil
}

func (d *Database) Delete(key string) (bool, error) {
	return false, nil
}

// New creates a new database
func New(dl logbook.Logbook, databaseName string) (*Database, error) {
	db := &Database{
		Dlog: dl,
	}
	_, err := os.Stat(config.DBFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(config.DBFolder, 0777)
			if err != nil {
				return nil, errors.New("error")
			}
		}
	}

	err = os.Chdir(config.DBFolder)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error")
	}

	f, err := os.Create(databaseName)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error")
	}

	defer f.Close()

	db.DatabaseName = databaseName
	db.Filename = databaseName
	db.Version = 001
	db.File = f
	db.Data = map[string][]byte{}

	wr := bufio.NewWriter(f)

	_, err = fmt.Fprintf(wr, "DolceDB.%d", config.DBVersion)
	if err != nil {
		fmt.Println(err)
	}
	wr.Flush()

	err = db.RebuildMap()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Delete db file
func Delete(db string) bool {
	err := os.Remove(db)
	if err != nil {
		return false
	}
	return true

}

func ListDBs() {

}

// RebuildMap is rebuilds the in memory map from the log
func (d *Database) RebuildMap() error {
	d.dbMutex.Lock()
	defer d.dbMutex.Unlock()

	logData, err := d.Dlog.GetState()
	if err != nil {
		log.Fatal("Could not get log.")
		return err
	}

	d.Data = logData

	return nil
}
