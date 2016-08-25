package logbook

type Index uint64

type Logbook interface {
	// Append adds a line to the log
	Append(string, []byte) error

	// GetFromIndex returns all the lines after a specific index
	GetFromIndex(Index) ([]string, error)

	// Replays all the log and gets a map of the current state
	GetState() (map[string][]byte, error)
}
