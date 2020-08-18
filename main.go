package main

import (
	"log"
	"os"
	"strconv"
	"time"

	internal_cache "github.com/sinuspower/image-previewer/internal/cache"
	internal_settings "github.com/sinuspower/image-previewer/internal/settings"
)

var (
	settings *internal_settings.Settings
	cache    internal_cache.Cache
)

func main() {
	settings = new(internal_settings.Settings)
	err := settings.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	cache, err = internal_cache.NewCache(settings.GetCacheSize(), "cache")
	if err != nil {
		log.Fatal("can not create cache:", err)
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
