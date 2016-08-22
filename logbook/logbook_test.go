package logbook

import "testing"

func TestNew(t *testing.T) {
	_, _, err := New("../data_test", "db_test.log")
	if err != nil {
		t.Error("Expected logbook, error ", err)
	}
}

func TestSet(t *testing.T) {
	dlog, found, err := New("../data_test", "db_test.log")
	if err != nil {
		t.Error("Expected logbook, error ", err)
	}

	if found == false {
		t.Error("Expected logbook file to exist, but not found")
	}

	err = dlog.Set("Key1", []byte("TestVal"))
	if err != nil {
		t.Error("Expect to set value, but got error: ", err)
	}
}
