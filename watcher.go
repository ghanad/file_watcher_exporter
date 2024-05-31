package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
)

func WatchPaths(config *Config) {
	for _, watcherConfig := range config.Watcher {
		initializeCurrentMetrics(watcherConfig)
		go watchPath(watcherConfig)
	}
}

func initializeCurrentMetrics(watcherConfig struct {
	Name   string              `yaml:"name"`
	Path   string              `yaml:"path"`
	Format string              `yaml:"format"`
	Labels []map[string]string `yaml:"labels"`
}) {
	var labels prometheus.Labels = make(map[string]string)
	labels["path"] = watcherConfig.Path
	for _, label := range watcherConfig.Labels {
		for k, v := range label {
			labels[k] = v
		}
	}
	updateCurrentMetrics(watcherConfig.Path, labels, watcherConfig.Format)
}

func watchPath(watcherConfig struct {
	Name   string              `yaml:"name"`
	Path   string              `yaml:"path"`
	Format string              `yaml:"format"`
	Labels []map[string]string `yaml:"labels"`
}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("File created: %s\n", event.Name)
					handleFileCreation(event.Name, watcherConfig)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("File deleted: %s\n", event.Name)
					handleFileDeletion(event.Name, watcherConfig)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	log.Printf("Watching path: %s\n", watcherConfig.Path)
	err = watcher.Add(watcherConfig.Path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func handleFileCreation(filePath string, watcherConfig struct {
	Name   string              `yaml:"name"`
	Path   string              `yaml:"path"`
	Format string              `yaml:"format"`
	Labels []map[string]string `yaml:"labels"`
}) {
	log.Printf("Handling file creation for: %s\n", filePath)
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	formats := strings.Split(watcherConfig.Format, ",")
	for _, format := range formats {
		format = strings.TrimSpace(format) // Remove any leading/trailing spaces
		log.Printf("Checking format: %s against extension: %s\n", format, ext)
		if ext == format {
			log.Printf("File format matched: %s\n", filePath)
			updateMetrics(filePath, watcherConfig)
			return
		}
	}
	log.Printf("File format did not match: %s\n", filePath)
}

func handleFileDeletion(filePath string, watcherConfig struct {
	Name   string              `yaml:"name"`
	Path   string              `yaml:"path"`
	Format string              `yaml:"format"`
	Labels []map[string]string `yaml:"labels"`
}) {
	log.Printf("Handling file deletion for: %s\n", filePath)
	var labels prometheus.Labels = make(map[string]string)
	labels["path"] = watcherConfig.Path
	for _, label := range watcherConfig.Labels {
		for k, v := range label {
			labels[k] = v
		}
	}
	updateCurrentMetrics(watcherConfig.Path, labels, watcherConfig.Format)
}
