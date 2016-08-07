package database

import "testing"

const dbName string = "TestDB"

func TestNew(t *testing.T) {
	database, err := New(dbName)
	if err != nil {
		t.Error("Expected database, got ", database)
	}
}
