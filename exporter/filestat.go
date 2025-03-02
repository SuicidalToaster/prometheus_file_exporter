package exporter

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var FileHash = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "pfe_file_hash",
}, []string{"hash", "path"})

func GetFileList(p []string) {
	for _, v := range p {
		go func() {
			for {
				GetFileHash(v)
				time.Sleep(5 * time.Second)
			}
		}()
	}
}

func GetFileHash(p string) {
	f, err := os.Open(p)
	if err != nil {
		FileHash.WithLabelValues("N/A", p).Set(0)
	}
	a, err := f.Stat()
	defer f.Close()
	if err != nil || a.IsDir() {
		FileHash.WithLabelValues("N/A", p).Set(0)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		FileHash.WithLabelValues("N/A", p).Set(0)
	}
	h := sha256.New()
	h.Write(b)
	s := base64.URLEncoding.EncodeToString(h.Sum(nil))
	FileHash.WithLabelValues(s, p).Set(1)
}
