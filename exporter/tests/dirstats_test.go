package tests

import (
	"fmt"
	"github.com/SuicidalToaster/prometheus_file_exporter/exporter"
	"testing"
	"time"
)

func TestGetCurrentDir(t *testing.T) {

	start := time.Now()
	// GetCurrentDir("/home/Dan/gitbucket/prometheus_file_exporter")
	count, err := exporter.GetTotalFiles("/")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetCurrentDir took %v with %d files\n", time.Since(start), count)
	// GetCurrentDir("/home/Dan/a")
}
