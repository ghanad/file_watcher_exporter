package main

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	configPath, showHelp, printConfig := parseArguments()

	if showHelp {
		printUsage()
		return
	}

	if printConfig {
		printSampleConfig()
		return
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up logging to file with rotation
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.LogPath,
		MaxSize:    10,   // megabytes after which new file is created
		MaxBackups: 7,    // number of backups
		MaxAge:     1,    // days
		Compress:   true, // compress the backups
	})

	log.Println("Starting application...")
	go StartMetricsServer(config.Exporter.Port, config.Exporter.Endpoint)
	WatchPaths(config)

	select {} // Block forever
}
