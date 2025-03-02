package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/SuicidalToaster/prometheus_file_exporter/exporter/v2"
)

func TestGetCurrentDir(t *testing.T) {

	start := time.Now()
	// GetCurrentDir("/home/Dan/gitbucket/prometheus_file_exporter")
	count, err := v2.GetTotalFiles("/")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GetCurrentDir took %v with %d files\n", time.Since(start), count)
	// GetCurrentDir("/home/Dan/a")
}
