package logbook

import "testing"

func TestNew(t *testing.T) {
	_, _, err := New("data_test", "db_test.log")
	if err != nil {
		t.Error("Expected logbook, error ", err)
	}
}
