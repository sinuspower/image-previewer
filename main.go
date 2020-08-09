package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	port, err := strconv.ParseUint(os.Getenv("IMAGE_PREVIEWER_PORT"), 10, 16)
	if err != nil {
		log.Fatalf("Error parsing IMAGE_PREVIEWER_PORT environment variable (%s)", err)
	}

	cacheSize, err := strconv.ParseUint(os.Getenv("IMAGE_PREVIEWER_CACHE_SIZE"), 10, 16)
	if err != nil {
		log.Fatalf("Error parsing IMAGE_PREVIEWER_CACHE_SIZE environment variable (%s)", err)
	}

	//-------------------------------------------
	/*
		logFile, err := os.Create("log.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := logFile.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	*/
	// log.SetOutput(logFile)
	//-------------------------------------------

	server := NewServer(uint16(port), uint16(cacheSize), os.Stdout)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
