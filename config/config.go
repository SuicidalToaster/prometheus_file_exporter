package config

import (
	"flag"
)

type ExporterConfig struct {
	Addr             string
	FilePaths        []string
	ExcludeFilePaths []string
	WalkDepth        int
}

type arrayFlags []string

var (
	addrFlag         = flag.String("port", "9111", "set listening port. Required")
	filePaths        arrayFlags
	excludeFilePaths arrayFlags
	walkDepth        = flag.Int("depth", -1, "--depth 1 --observe /data will observe only files in /data dir and will not go into dirs under /data. Value of -1 will walk through all children directories")
)

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func GetConfig() ExporterConfig {
	var cfg ExporterConfig
	flag.Var(&filePaths, "observe", "set watched dirs \n --observe ./ --observe=/data")
	flag.Var(&excludeFilePaths, "exclude", "exclude watched dirs \n --observe /data --exclude=/data/bad-data")
	flag.Parse()
	cfg.Addr = *addrFlag
	cfg.ExcludeFilePaths = excludeFilePaths
	cfg.FilePaths = filePaths
	cfg.WalkDepth = *walkDepth
	return cfg
}
