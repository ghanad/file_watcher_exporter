package main

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

func printUsage() {
	fmt.Println("Usage: FileMetricsWatcher [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("  -h                  Show help message")
	fmt.Println("  --config.file PATH  Path to the configuration file")
	fmt.Println("  --print-config      Print a sample YAML configuration")
}

func printSampleConfig() {
	sampleConfig := Config{
		Exporter: struct {
			Port     int    `yaml:"port"`
			Endpoint string `yaml:"endpoint"`
		}{
			Port:     2012,
			Endpoint: "/metric",
		},
		LogPath: "path/to/your/logfile.log",
		Watcher: []struct {
			Name   string              `yaml:"name"`
			Path   string              `yaml:"path"`
			Format string              `yaml:"format"`
			Labels []map[string]string `yaml:"labels"`
		}{
			{
				Name:   "test",
				Path:   "/mnt/extra/dir1",
				Format: "all",
				Labels: []map[string]string{
					{"path": "test"},
					{"name": "test"},
				},
			},
			{
				Name:   "test2",
				Path:   "/mnt/extra/dir2",
				Format: "py, csv",
				Labels: []map[string]string{
					{"project": "myProj"},
					{"app_name": "app1"},
				},
			},
		},
	}

	data, err := yaml.Marshal(&sampleConfig)
	if err != nil {
		log.Fatalf("Error generating sample config: %v", err)
	}
	fmt.Println(string(data))
}

func parseArguments() (string, bool, bool) {
	var configPath string
	flag.StringVar(&configPath, "config.file", "config.yaml", "Path to the configuration file")
	showHelp := flag.Bool("h", false, "Show help message")
	printConfig := flag.Bool("print-config", false, "Print a sample YAML configuration")

	flag.Parse()

	return configPath, *showHelp, *printConfig
}
