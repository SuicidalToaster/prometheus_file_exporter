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
	switch cfg.ShowRootOnly {
	case true:
		for _, v := range cfg.FilePaths {
			go func() {
				for {
					var cnt int
					err := filepath.WalkDir(v, func(path string, d fs.DirEntry, err error) error {
						if !d.IsDir() {
							cnt++
						}
						return nil
					})
					if err != nil {
						log.Println(err)
					}
					PathFileCount.WithLabelValues(v).Set(0)
					PathFileCount.WithLabelValues(v).Set(float64(cnt))
					cnt = 0
					time.Sleep(10 * time.Second)
				}
			}()
		}
	}
}
