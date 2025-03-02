package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
	v2 "github.com/SuicidalToaster/prometheus_file_exporter/exporter/v2"
)

func main() {
	// runtime.GOMAXPROCS(1)
	// debug.SetGCPercent(20)
	// debug.SetMemoryLimit(1024 * 1024 * 1024 * 1024)
	conf := config.GetConfig()
	mux := http.NewServeMux()
	for _, v := range conf.FilePaths {
		go func() {
			for {
				start := time.Now()
				v2.GetTotalFiles(v)
				fmt.Println(time.Since(start))
			}
		}()
	}
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
