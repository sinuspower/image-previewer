package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

var settings *Settings

func main() {
	settings = new(Settings)
	err := settings.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().Format("2006-01-02_15:04:05")
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0700); err != nil {
			log.Fatal(err)
		}
	}
	logFile, err := os.Create("logs/" + now + "_" +
		strconv.Itoa(settings.GetPort()) + "_image-previewer_log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	log.SetOutput(logFile)

	server := NewServer(settings.GetPort(), settings.GetCacheSize(), logFile)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
