package main

import (
	"context"
	"fmt"
	"github.com/SuicidalToaster/prometheus_file_exporter/exporter"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/SuicidalToaster/prometheus_file_exporter/config"
)

var cfg = config.InitConfig()

func main() {
	//runtime.GOMAXPROCS(1)
	//debug.SetGCPercent(20)
	//debug.SetMemoryLimit(1024 * 1024 * 1024 * 1024 * 3)
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
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
	go LaunchFileCount()
	go LaunchFileHash()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "prometheus_file_exporter. Exports various fs metrics")
	})
	srv := http.Server{
		Addr: ":" + cfg.GetString("server.port"),
		// ErrorLog: log.Default(),
		Handler: mux,
	}
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	select {
	case <-ctx.Done():
		err = srv.Shutdown(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func LaunchFileCount() {
	wg := sync.WaitGroup{}
	for {
		for _, v := range cfg.GetStringSlice("DirPaths") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				t := time.Now()
				c, err := exporter.GetTotalFiles(v)
				if err != nil {
					log.Println(err)
				}
				log.Println(c, v, time.Since(t).Seconds())

			}()
		}
		wg.Wait()
	}
}

func LaunchFileHash() {
	wg := sync.WaitGroup{}
	for {
		for _, v := range cfg.GetStringSlice("HashFiles") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				t := time.Now()
				exporter.GetFileHash(v)
				log.Println("Hash", v, time.Since(t).Seconds())
			}()
		}
		wg.Wait()
	}
}
