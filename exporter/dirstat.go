package exporter

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
)

var PathFileCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "path_file_count",
	Help: "Shows cumulative directory file count like du",
}, []string{"path", "exactCount"})

func GetFSMetrics(cfg config.ExporterConfig) {

	for _, v := range cfg.FilePaths {
		go func() {
			for {
				start := time.Now()
				var fileCount int
				err := filepath.WalkDir(v, func(path string, d fs.DirEntry, err error) error {
					switch d.IsDir() {
					case true:
						if slices.Contains(cfg.ExcludeFilePaths, path) {
							return filepath.SkipDir
						}
					case false:
						fileCount++
					}

					return nil
				})
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("%s\n", time.Since(start))
				PathFileCount.Reset()
				PathFileCount.WithLabelValues(v, strconv.Itoa(fileCount)).Set(float64(fileCount))
				fileCount = 0
				time.Sleep(5 * time.Second)
			}
		}()
	}

}
