package database

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/superdecimal/dolce/logbook"
)

const dbName string = "TestDB"

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dLog := logbook.NewMockLogbook(ctrl)
	dLog.EXPECT().GetAll()

	database, err := New(dLog, dbName)

	if err != nil {
		t.Error("Expected database, got ", database)
	} else {
		fmt.Println("Success")
	}
}
