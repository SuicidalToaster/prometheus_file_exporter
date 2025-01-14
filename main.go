package main

import (
	"fmt"
	"net/http"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
	"github.com/SuicidalToaster/prometheus_file_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	conf := config.GetConfig()
	go exporter.GetFSMetrics(conf)
	go exporter.GetFileList(conf.HashFiles)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "prometheus_file_exporter. Exports various fs metrics")
	})
	srv := http.Server{
		Addr: ":" + conf.Addr,
		// ErrorLog: log.Default(),
		Handler: mux,
	}
	srv.ListenAndServe()
}
