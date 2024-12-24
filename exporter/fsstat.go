package exporter

import (
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
)

var PathFileCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "path_file_count",
	Help: "Shows cumulative directory file count like du",
}, []string{"path"})

func GetFSMetrics(cfg config.ExporterConfig) {

	for _, v := range cfg.FilePaths {
		go func() {
			for {
				var fileCount int
				err := filepath.WalkDir(v, func(path string, d fs.DirEntry, err error) error {
					switch d.IsDir() {
					case true:

					case false:
						fileCount++
					}
					return nil
				})
				if err != nil {
					log.Println(err)
				}
				PathFileCount.WithLabelValues(v).Set(0)
				PathFileCount.WithLabelValues(v).Set(float64(fileCount))
				fileCount = 0
				time.Sleep(10 * time.Second)
			}
		}()
	}

}
