package exporter

import (
	"github.com/SuicidalToaster/prometheus_file_exporter/config"
)

func Collector(cfg config.ExporterConfig) {
	go GetFSMetrics(cfg)
}
