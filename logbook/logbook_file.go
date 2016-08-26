// Package logbook provides an append only log and functions to handle transactions with it.
package logbook

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrNotFound     = errors.New("Error log file not found or created.")
	ErrFolderCreate = errors.New("Error creating log folder")
	ErrFileCreate   = errors.New("Error creating log file.")
	ErrReadFile     = errors.New("Error reading the log file.")
	ErrWriteFile    = errors.New("Error writing to the log file")
)

// LogbookFile is the basic structure used by the log
// - filename is the log filename
// - path is the folder path of the log
// - file is a point to the log file
// - index is the last number used as an index in the log
// - logMutex is a mutex to lock/unlock writing to the log
type LogbookFile struct {
	filename string
	path     string
	file     *os.File
	version  int
	index    Index
	logMutex sync.Mutex
}

func ToIndex(i string) (Index, error) {
	temp, _ := strconv.ParseUint(i, 0, 64)
	return Index(temp), nil
}

//New creates a new log
func New(fp, fn string) (*LogbookFile, bool, error) {
	var i Index

	dlog := &LogbookFile{
		version:  1,
		filename: fn,
		path:     fp,
		index:    0,
	}

	var filepath = dlog.path + "/" + dlog.filename

	_, err := os.Stat(dlog.path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(dlog.path, 0777)
			if err != nil {
				return nil, false, ErrFolderCreate
			}
		}
	}

	//Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// If not create it
		f, err := os.Create(filepath)
		if err != nil {
			return nil, false, ErrFileCreate
		}

		dlog.file = f

		fmt.Println("Log file not found and created.")

		return dlog, false, nil
	}

	// If file exists open it
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, true, ErrNotFound
	}

	//Check logfile index
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
	}

	dlog.index = i
	dlog.file = f

	fmt.Println("Log file found.")

	return dlog, true, nil
}

// Set appened to the log
func (l *LogbookFile) Append(key string, value []byte) error {
	l.logMutex.Lock()
	defer l.logMutex.Unlock()

	wr := bufio.NewWriter(l.file)
	_, err := fmt.Fprintf(wr, "%d  S %q %q\n", l.index, key, value)
	if err != nil {
		fmt.Println(err)
	}

	l.index++

	err = wr.Flush()
	if err != nil {
		return ErrWriteFile
	}

	return nil
}

// GetFromIndex returns the log after a specific index.
// TODO implement a better file parser
func (l *LogbookFile) GetFromIndex(index Index) (<-chan string, error) {
	l.logMutex.Lock()
	out := make(chan string)
	var i Index

	f, err := os.Open(l.filename)
	if err != nil {
		return nil, ErrReadFile
	}

	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return nil, ErrReadFile
	}

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			sepIndex := strings.Index(line, " ")
			id, _ := ToIndex(line[:sepIndex])

			if index <= id {
				out <- scanner.Text()
			}
			i++
		}

		close(out)
		f.Close()
		l.logMutex.Unlock()
	}()

	return out, nil
}

// GetAll returns a map with the latest state
func (l *LogbookFile) GetAll() (<-chan string, error) {
	l.logMutex.Lock()
	out := make(chan string)

	var i Index
	i = 0

	f, err := os.Open(l.filename)
	if err != nil {
		return nil, ErrReadFile
	}

	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return nil, ErrReadFile
	}

	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
			i++
		}

		l.index = i

		fmt.Println("Index: ", i)

		close(out)
		f.Close()
		l.logMutex.Unlock()
	}()

	return out, nil
}
