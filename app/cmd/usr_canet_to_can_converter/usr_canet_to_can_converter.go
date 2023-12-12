package main

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	usrcanettocan "github.com/dmnunez1993/usr-canet200-to-can-converter"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {

	} else {
		log.Info("Loading .env")
		godotenv.Load()
	}

	converter := usrcanettocan.NewUsrCanetConverter()
	go converter.Run()

	go usrcanettocan.ServeRestApi()

	usrcanettocan.LoopForever()
}
