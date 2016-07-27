package storage

import "testing"

const dbName string = "TestDB"

func TestCreateDBFile(t *testing.T) {
	database, err := CreateDBFile(dbName)
	if err != nil {
		t.Error("Expected database, got ", database)
	}
}
