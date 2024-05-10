package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	ClearCache()
	CheckPathDb()
	db, err := LoadDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Initializing...")

	StartConnection(os.Getenv("PHONE_NUMBER"), db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Exiting...")
}
