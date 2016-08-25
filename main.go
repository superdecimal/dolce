package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/superdecimal/dolce/config"
	"github.com/superdecimal/dolce/database"
	"github.com/superdecimal/dolce/logbook"
	"github.com/superdecimal/dolce/networking"
)

func main() {
	dlog, _, err := logbook.New("data", "db.log")
	if err != nil {
		log.Fatal("No log could be created or found")
	}

	db, err := database.New(dlog, config.DBFilename)
	if err != nil {
		log.Fatal("Could not create db file")
		return
	}

	// for i := 0; i < 1000; i++ {
	// 	db.Set("TestKey6", fmt.Sprintf("TestValue %d", i))
	// }

	data, _ := db.Read("TestKey6")
	fmt.Println(data)

	go networking.StartServer()
	go networking.StartTCPServer()

	c := make(chan os.Signal, 1000)
	signal.Notify(c, os.Interrupt)

	func() {
		for _ = range c {
			fmt.Println("Exiting...")
			return
		}
	}()
}
