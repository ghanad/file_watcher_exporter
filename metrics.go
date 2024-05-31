package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fileCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "path_watcher_total",
			Help: "Total number of files watched",
		},
		[]string{"path", "project", "app_name"},
	)
	fileSize = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "path_watcher_size_total",
			Help: "Total size of files watched",
		},
		[]string{"path", "project", "app_name"},
	)
	currentFileCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "path_watcher_current_file_count",
			Help: "Current number of files in the watched directories",
		},
		[]string{"path", "project", "app_name"},
	)
	currentFileSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "path_watcher_current_file_size",
			Help: "Current total size of files in the watched directories",
		},
		[]string{"path", "project", "app_name"},
	)
)

func init() {
	log.Println("Registering metrics...")
	prometheus.MustRegister(fileCount)
	prometheus.MustRegister(fileSize)
	prometheus.MustRegister(currentFileCount)
	prometheus.MustRegister(currentFileSize)
	log.Println("Metrics registered.")
}

func updateMetrics(filePath string, watcherConfig struct {
	Name   string              `yaml:"name"`
	Path   string              `yaml:"path"`
	Format string              `yaml:"format"`
	Labels []map[string]string `yaml:"labels"`
}) {
	log.Printf("Updating metrics for file: %s\n", filePath)
	var labels prometheus.Labels = make(map[string]string)
	labels["path"] = watcherConfig.Path
	for _, label := range watcherConfig.Labels {
		for k, v := range label {
			labels[k] = v
		}
	}
	log.Printf("Labels for metrics: %v\n", labels)
	fileCount.With(labels).Inc()
	log.Printf("Incremented file count metric for path: %s, labels: %v\n", watcherConfig.Path, labels)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Println("Error getting file info:", err)
		return
	}
	fileSize.With(labels).Add(float64(fileInfo.Size()))
	log.Printf("Incremented file size metric for path: %s, labels: %v, size: %d\n", watcherConfig.Path, labels, fileInfo.Size())

	updateCurrentMetrics(watcherConfig.Path, labels, watcherConfig.Format)
}

func updateCurrentMetrics(dirPath string, labels prometheus.Labels, formats string) {
	totalSize := int64(0)
	fileCount := 0
	formatSet := make(map[string]struct{})
	for _, format := range strings.Split(formats, ",") {
		formatSet[strings.TrimSpace(format)] = struct{}{}
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if _, ok := formatSet[ext]; ok {
				fileCount++
				totalSize += info.Size()
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Error walking the path:", err)
		return
	}
	currentFileCount.With(labels).Set(float64(fileCount))
	currentFileSize.With(labels).Set(float64(totalSize))
	log.Printf("Updated current file count to %d and size to %d for path: %s, labels: %v\n", fileCount, totalSize, dirPath, labels)
}

func StartMetricsServer(port int, endpoint string) {
	if endpoint == "" {
		log.Fatal("Endpoint pattern must not be empty")
	}

	http.Handle(endpoint, promhttp.Handler())
	log.Printf("Starting metrics server on port %d at endpoint %s\n", port, endpoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
