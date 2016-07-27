package dolcelog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

/*
Basic structure to be used for the log
- filename is the log filename
- path is the folder path of the log
- file is a point to the log file
- index is the last number used as an index in the log
- logMutex is a mutex to lock/unlock writing to the log
*/
type DolceLog struct {
	filename string
	path     string
	file     *os.File
	version  int
	index    uint64
}

var dlog DolceLog

func init() {
	dlog.version = 1
	dlog.filename = "db.log"
	dlog.path = "data"
	// TODO change this reflect the index of the log when the server is restarted
	dlog.index = 0

	var filepath = dlog.path + "/" + dlog.filename

	//Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// If not create it
		f, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Log file not found and created")
		dlog.file = f
		return
	}

	// If file exists open it
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Log file found")
	dlog.file = f

	//TODO ReadFile and return map and index
}

func (l *DolceLog) Set(key string, value []byte) {
	wr := bufio.NewWriter(l.file)

	_, err := fmt.Fprintf(wr, "%d  S %s %s\n", l.index, key, value)
	if err != nil {
		fmt.Println(err)
	}

	l.index++

	err = wr.Flush()
	if err != nil {
		log.Fatal("Data not writen in file.")
	}
}

func GetLogInst() *DolceLog {
	return &dlog
}

func (l *DolceLog) RecreateMap() {

}
