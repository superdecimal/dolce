package logbook

type Index uint64

type Logbook interface {
	Set(string, []byte) error
	GetFromIndex(Index) ([]string, error)
	GetAll() ([]string, error)
}
