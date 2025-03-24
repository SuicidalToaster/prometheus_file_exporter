package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var PathFileCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "pfe_path_file_count",
	Help: "Shows cumulative directory file count like du",
}, []string{"path"})

type Task struct {
	ID int
}

func (wp *Task) Do(f func()) {
	f()
}

type WorkerPool struct {
	Tasks       []Task
	concurrency int
	taskChan    chan Task
	wg          sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for task := range wp.taskChan {
		task.Do()
	}
}

func (wp *WorkerPool) Start() {
	wp.taskChan = make(chan Task, len(wp.Tasks))
	for i := range wp.concurrency {
		go wp.worker()
	}

}

func GetTotalFiles(path string) (int, error) {
	res := make(chan int)
	total := 0
	fi, err := os.Stat(path)
	if err != nil {
		// log.Printf("Error getting current directory stats: %v", err)
		return 0, err
	}
	if fi.IsDir() {
		go startCount(path, res)

	} else {
		// log.Printf("Not a directory: %s. Is %s\n", path, fi.Mode().Type().String())
		close(res)
		return 0, fmt.Errorf("%s is not a directory", path)
	}
	for i := range res {
		total += i
	}
	PathFileCount.WithLabelValues(path).Set(float64(total))
	return total, err
}

func startCount(path string, ch chan int) {
	f, err := os.Stat(path)
	if err != nil {
		// log.Printf("Error opening directory %s: %v", path, err)
		close(ch)
		return
	}
	switch f.IsDir() {
	case true:
		countFiles(path, ch)
	case false:
		close(ch)
		return
	}
	close(ch)
}

var lock sync.Mutex

func countFiles(path string, ch chan int) {

	var wg = sync.WaitGroup{}
	var totalInDir int
	de, err := os.ReadDir(path)
	if err != nil {
		// log.Printf("Error reading directory %s: %v", path, err)
		return
	}
	wg.Add(len(de))
	go func() {
		lock.Lock()
		for _, v := range de {
			switch v.Type() {
			case os.ModeDir:
				go func() {
					countFiles(filepath.Join(path, v.Name()), ch)
					wg.Done()
				}()
				// break
			case 0:
				totalInDir++
				wg.Done()
				// break
			default:
				wg.Done()
				// break
			}
		}
		lock.Unlock()
	}()

	wg.Wait()

	ch <- totalInDir
}
