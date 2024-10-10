package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	loader "tomoribot-geminiai-version/src/main"
	"tomoribot-geminiai-version/src/utils"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Initializing...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	utils.ClearCache()
	utils.CheckPathDb()

	loader.StartConnection(os.Getenv("PHONE_NUMBER"), nil)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Exiting...")
}
