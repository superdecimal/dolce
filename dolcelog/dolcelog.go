// Package dolcelog provides an append only log and functions to handle transactions with it.
package dolcelog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
DolceLog is the basic structure used by the log
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
	logMutex sync.Mutex
}

var dlog DolceLog

// Initialiazes the log
func init() {
	dlog.version = 1
	dlog.filename = "db.log"
	dlog.path = "data"
	// TODO change this to reflect the index of the log when the server is restarted
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

// Set appned to the log
func (l *DolceLog) Set(key string, value []byte) {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

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

// GetLogInst returns the log instance.
func GetLogInst() *DolceLog {
	return &dlog
}

// GetFromIndex returns the log after a specific index.
// TODO implement a better file parser
func (l *DolceLog) GetFromIndex(index int) ([]string, error) {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	var result = make([]string, 0)
	i := 0

	f, err := os.Open(dlog.filename)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		sepIndex := strings.Index(line, " ")
		id, err := strconv.Atoi(line[:sepIndex])
		if err != nil {
			log.Fatal("Index retrieval error.")
		}

		if index <= id {
			result = append(result, line)
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result, nil
}

// GetAll returns a slice of the with the whole log.
func (l *DolceLog) GetAll() ([]string, error) {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	var result = make([]string, 0)
	i := 0

	f, err := os.Open(dlog.filename)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		result = append(result, scanner.Text())
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result, nil
}
