package main

import (
	"fmt"
	v2 "github.com/SuicidalToaster/prometheus_file_exporter/exporter/v2"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
)

var cfg = config.InitConfig()

func main() {
	//runtime.GOMAXPROCS(1)
	//debug.SetGCPercent(20)
	//debug.SetMemoryLimit(1024 * 1024 * 1024 * 1024)

	go func() {
		for {
			time.Sleep(time.Second * 5)
			cfg.WatchConfig()
			cfg.OnConfigChange(func(e fsnotify.Event) {
				log.Println("config file changed", e.Name)
			})
		}
	}()
	mux := http.NewServeMux()
	go func() {
		LaunchFileCount()
	}()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "prometheus_file_exporter. Exports various fs metrics")
	})
	srv := http.Server{
		Addr: ":" + cfg.GetString("server.port"),
		// ErrorLog: log.Default(),
		Handler: mux,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func LaunchFileCount() {
	for {
		for _, v := range cfg.GetStringSlice("DirPaths") {
			t := time.Now()
			c, err := v2.GetTotalFiles(v)
			if err != nil {
				log.Println(err)
			}
			println(c, v, time.Since(t).Seconds())
		}
		time.Sleep(25 * time.Second)
	}
}
