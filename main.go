package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
	"github.com/SuicidalToaster/prometheus_file_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var conf = config.GetConfig()

func main() {
	start := time.Now()
	go exporter.GetFSMetrics(conf)
	fmt.Printf("%s", time.Since(start))
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
