package exporter

import (
	"log"
	"os"
	"slices"
	"strings"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var PathFileCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "path_file_count",
	Help: "Shows cumulative directory file count like du",
}, []string{"path"})

func GetFSMetrics(cfg config.ExporterConfig) {
	for _, v := range cfg.FilePaths {
		go func() {
			PathFileCount.WithLabelValues(v).Set(0)
			CountFiles(v, v, &cfg)
		}()
	}
}

func CountFiles(rootDir string, curDir string, cfg *config.ExporterConfig) {

	if !strings.Contains(curDir, rootDir) {
		curDir = rootDir + strings.TrimSuffix(curDir, "/")
	}
	// absDir = absDir + "/" + relPath
	tree, err := os.ReadDir(curDir)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	if cfg.WalkDepth == 0 {
		return
	}
	cfg.WalkDepth = cfg.WalkDepth - 1
	for _, v := range tree {
		if v.IsDir() {
			if slices.Contains(cfg.ExcludeFilePaths, curDir+"/"+v.Name()) {
				continue
			}
			go CountFiles(rootDir, curDir+"/"+v.Name(), cfg)
		} else {
			switch curDir == rootDir {
			case true:
				PathFileCount.WithLabelValues(curDir).Inc()
			case false:
				PathFileCount.WithLabelValues(curDir).Inc()
				PathFileCount.WithLabelValues(rootDir).Inc()
			}
		}
	}
}
