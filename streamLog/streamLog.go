package streamlog

import (
	"os"

	"github.com/superdecimal/dolce/config"
)

type Log struct {
	filename string
	path     string
	file     *os.File
	version  int
}

func init() {
	var log Log

	log.version = 1
	log.filename = config.LogFilename
	log.path = "data"
	if _, err := os.Stat(log.path + "/" + log.filename); os.IsNotExist(err) {
		_, err := os.Create(log.path + "/" + log.filename)
		if err != nil {

		}
	}
}

func (l *Log) Set(key string, value []byte) bool {
	return true
}
