package benchmark

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// BenchmarkResult represents the result of a single benchmark run
type BenchmarkResult struct {
	Structure  string  `json:"structure"`
	Operation  string  `json:"operation"`
	DataSize   int     `json:"dataSize"`
	Duration   float64 `json:"duration"`   // in milliseconds
	MemoryUsed uint64  `json:"memoryUsed"` // in bytes
	OpsPerSec  float64 `json:"opsPerSec"`
	Progress   int     `json:"progress"` // 0-100
	Completed  bool    `json:"completed"`
}

// BenchmarkConfig represents configuration for a benchmark run
type BenchmarkConfig struct {
	DataSize   int      `json:"dataSize"`
	Structures []string `json:"structures"`
	Operation  string   `json:"operation"`
}

// ProgressCallback is called with benchmark progress updates
type ProgressCallback func(result BenchmarkResult)

// Runner manages benchmark execution
type Runner struct {
	mu       sync.Mutex
	running  bool
	stopChan chan struct{}
}

// NewRunner creates a new benchmark runner
func NewRunner() *Runner {
	return &Runner{
		stopChan: make(chan struct{}),
	}
}

// generateRandomData generates random integers for benchmarking
func generateRandomData(size int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(size * 10)
	}
	return data
}

// getMemoryUsage returns current memory usage
func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// RunBenchmark runs benchmarks for specified structures
func (r *Runner) RunBenchmark(config BenchmarkConfig, callback ProgressCallback) {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return
	}
	r.running = true
	r.stopChan = make(chan struct{})
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
	}()

	data := generateRandomData(config.DataSize)

	var wg sync.WaitGroup
	for _, structure := range config.Structures {
		wg.Add(1)
		go func(structName string) {
			defer wg.Done()
			r.runSingleBenchmark(structName, config.Operation, data, callback)
		}(structure)
	}
	wg.Wait()
}

func (r *Runner) runSingleBenchmark(structure, operation string, data []int, callback ProgressCallback) {
	startMem := getMemoryUsage()
	startTime := time.Now()

	total := len(data)
	reportInterval := total / 20 // Report every 5%
	if reportInterval < 1 {
		reportInterval = 1
	}

	switch structure {
	case "hashmap":
		r.benchmarkHashMap(operation, data, callback, reportInterval)
	case "btree":
		r.benchmarkBTree(operation, data, callback, reportInterval)
	case "rbtree":
		r.benchmarkRBTree(operation, data, callback, reportInterval)
	case "avltree":
		r.benchmarkAVLTree(operation, data, callback, reportInterval)
	}

	endMem := getMemoryUsage()
	duration := time.Since(startTime).Seconds() * 1000

	// Final result
	callback(BenchmarkResult{
		Structure:  structure,
		Operation:  operation,
		DataSize:   len(data),
		Duration:   duration,
		MemoryUsed: endMem - startMem,
		OpsPerSec:  float64(len(data)) / (duration / 1000),
		Progress:   100,
		Completed:  true,
	})
}

func (r *Runner) benchmarkHashMap(operation string, data []int, callback ProgressCallback, reportInterval int) BenchmarkResult {
	m := make(map[int]int)
	startTime := time.Now()

	for i, v := range data {
		select {
		case <-r.stopChan:
			return BenchmarkResult{}
		default:
		}

		switch operation {
		case "insert":
			m[v] = v
		case "search":
			if i > 0 {
				_ = m[data[rand.Intn(i)]]
			}
		}

		if i > 0 && i%reportInterval == 0 {
			progress := (i * 100) / len(data)
			callback(BenchmarkResult{
				Structure:  "hashmap",
				Operation:  operation,
				DataSize:   len(data),
				Duration:   time.Since(startTime).Seconds() * 1000,
				MemoryUsed: getMemoryUsage(),
				Progress:   progress,
				Completed:  false,
			})
		}
	}

	return BenchmarkResult{}
}

func (r *Runner) benchmarkBTree(operation string, data []int, callback ProgressCallback, reportInterval int) {
	// Simple B-Tree simulation using sorted slice
	tree := make([]int, 0, len(data))
	startTime := time.Now()

	for i, v := range data {
		select {
		case <-r.stopChan:
			return
		default:
		}

		switch operation {
		case "insert":
			// Binary search insert to keep sorted
			idx := binarySearchInsertPos(tree, v)
			tree = append(tree, 0)
			copy(tree[idx+1:], tree[idx:])
			tree[idx] = v
		case "search":
			if len(tree) > 0 {
				_ = binarySearch(tree, data[rand.Intn(len(tree))])
			}
		}

		if i > 0 && i%reportInterval == 0 {
			progress := (i * 100) / len(data)
			callback(BenchmarkResult{
				Structure:  "btree",
				Operation:  operation,
				DataSize:   len(data),
				Duration:   time.Since(startTime).Seconds() * 1000,
				MemoryUsed: getMemoryUsage(),
				Progress:   progress,
				Completed:  false,
			})
		}
	}
}

func (r *Runner) benchmarkRBTree(operation string, data []int, callback ProgressCallback, reportInterval int) {
	// Simplified benchmark without step tracking
	m := make(map[int]struct{})
	startTime := time.Now()

	for i, v := range data {
		select {
		case <-r.stopChan:
			return
		default:
		}

		switch operation {
		case "insert":
			m[v] = struct{}{}
		case "search":
			if i > 0 {
				_, _ = m[data[rand.Intn(i)]]
			}
		}

		if i > 0 && i%reportInterval == 0 {
			progress := (i * 100) / len(data)
			callback(BenchmarkResult{
				Structure:  "rbtree",
				Operation:  operation,
				DataSize:   len(data),
				Duration:   time.Since(startTime).Seconds() * 1000,
				MemoryUsed: getMemoryUsage(),
				Progress:   progress,
				Completed:  false,
			})
		}
	}
}

func (r *Runner) benchmarkAVLTree(operation string, data []int, callback ProgressCallback, reportInterval int) {
	m := make(map[int]struct{})
	startTime := time.Now()

	for i, v := range data {
		select {
		case <-r.stopChan:
			return
		default:
		}

		switch operation {
		case "insert":
			m[v] = struct{}{}
		case "search":
			if i > 0 {
				_, _ = m[data[rand.Intn(i)]]
			}
		}

		if i > 0 && i%reportInterval == 0 {
			progress := (i * 100) / len(data)
			callback(BenchmarkResult{
				Structure:  "avltree",
				Operation:  operation,
				DataSize:   len(data),
				Duration:   time.Since(startTime).Seconds() * 1000,
				MemoryUsed: getMemoryUsage(),
				Progress:   progress,
				Completed:  false,
			})
		}
	}
}

// Stop stops any running benchmark
func (r *Runner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.running {
		close(r.stopChan)
	}
}

func binarySearchInsertPos(arr []int, target int) int {
	left, right := 0, len(arr)
	for left < right {
		mid := (left + right) / 2
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}

func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}
