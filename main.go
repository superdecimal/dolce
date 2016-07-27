package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/superdecimal/dolce/config"
	"github.com/superdecimal/dolce/networking"
	"github.com/superdecimal/dolce/storage"
)

func main() {
	db, err := storage.CreateDBFile(config.DBFilename)
	if err != nil {
		log.Fatal("Could not create db file")
		return
	}

	db.Set("TestKey6", "TestValue100")

	data, _ := db.Read("TestKey6")
	fmt.Println(data)

	go networking.StartServer()
	go networking.StartTCPServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	func() {
		for _ = range c {
			fmt.Println("Exiting...")
			return
		}
	}()
}
