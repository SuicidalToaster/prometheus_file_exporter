package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
	"github.com/SuicidalToaster/prometheus_file_exporter/exporter"
)

var conf = config.GetConfig()

func main() {

	go exporter.GetFSMetrics(conf)
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
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
