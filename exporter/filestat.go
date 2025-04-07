package exporter

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io"
	"os"
)

var FileHash = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "pfe_file_hash",
}, []string{"hash", "path"})

func GetFileHash(p string) {
	f, err := os.Open(p)
	if err != nil {
		FileHash.WithLabelValues("N/A", p).Set(0)
		return
	}
	_, err = f.Stat()
	defer f.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			FileHash.WithLabelValues("N/A", p).Set(0)
		}
		return
	}
	b, err := io.ReadAll(f)
	if err != nil {
		FileHash.WithLabelValues("N/A", p).Set(0)
		return
	}
	h := sha256.New()
	h.Write(b)
	s := base64.URLEncoding.EncodeToString(h.Sum(nil))
	FileHash.DeleteLabelValues(s, p)
	FileHash.WithLabelValues(s, p).Set(1)
}
