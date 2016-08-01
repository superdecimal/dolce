package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/superdecimal/dolce/config"
	"github.com/superdecimal/dolce/dolcelog"
)

type Database struct {
	DatabaseName string
	Filename     string
	Path         string
	File         *os.File
	Version      int
	Data         map[string][]byte
	dbMutex      sync.Mutex
}

func init() {
}

func (d *Database) Set(key string, value string) error {
	d.dbMutex.Lock()
	defer d.dbMutex.Unlock()

	dlog := dolcelog.GetLogInst()
	data := []byte(value)
	dlog.Set(key, data)
	d.Data[key] = data

	return nil
}

func (d *Database) Read(key string) (string, error) {
	return string(d.Data[key]), nil
}

func (d *Database) Delete(key string) (bool, error) {
	return false, nil
}

// CreateDBFile creates the db file and returns a pointer to it
func CreateDBFile(databaseName string) (*Database, error) {
	var db Database
	_, err := os.Stat(config.DBFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(config.DBFolder, 0777)
			if err != nil {
				fmt.Println(err)
				return nil, errors.New("error")
			}
		}
	}

	err = os.Chdir(config.DBFolder)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(databaseName)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	db.DatabaseName = databaseName
	db.Filename = databaseName
	db.Version = 001
	db.File = f
	db.Data = make(map[string][]byte, 1000)

	wr := bufio.NewWriter(f)

	_, err = fmt.Fprintf(wr, "DolceDB.%d", config.DBVersion)
	if err != nil {
		fmt.Println(err)
	}
	wr.Flush()

	return &db, nil
}
func DeleteDBFile(db string) bool {
	err := os.Remove(db)
	if err != nil {
		return false
	}
	return true

}

func ListDBs() {

}

func RebuildMap() {

}
