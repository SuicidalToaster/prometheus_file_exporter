package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	v2 "github.com/SuicidalToaster/prometheus_file_exporter/exporter/v2"
)

func BenchmarkDirstatsV2(b *testing.B) {
	eg := errgroup.Group{}
	eg.SetLimit(10)

	paths := []string{"/"}
	wg := sync.WaitGroup{}
	wg.Add(len(paths))
	mymap := make(map[string]int)
	timemap := make(map[string]time.Duration)
	for _, path := range paths {
		go func() {
			start := time.Now()
			mymap[path], _ = v2.GetTotalFiles(path)
			timemap[path] = time.Since(start)
			wg.Done()
		}()
	}
	wg.Wait()
	for _, path := range paths {
		fmt.Printf("Total files in %s is %d. Done in %f seconds\n", path, mymap[path], timemap[path].Seconds())
	}
}
