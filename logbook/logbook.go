package logbook

type Index uint64

type Logbook interface {
	// Append adds a line to the log
	Append(string, []byte) error

	// GetFromIndex returns all the lines after a specific index
	GetFromIndex(Index) (<-chan string, error)

	// GetAll returns a channel and send the logbook line by line
	GetAll() (<-chan string, error)
}
