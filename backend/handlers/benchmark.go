package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"gin/benchmark"

	"github.com/gin-gonic/gin"
)

// BenchmarkRequest represents a request to start a benchmark
type BenchmarkRequest struct {
	DataSize   int      `json:"dataSize" binding:"required"`
	Structures []string `json:"structures" binding:"required"`
	Operation  string   `json:"operation" binding:"required"`
}

var (
	benchmarkRunner = benchmark.NewRunner()
	runnerMutex     sync.Mutex
)

// HandleBenchmarkSSE handles SSE connections for benchmark progress
func HandleBenchmarkSSE(c *gin.Context) {
	var req BenchmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Channel for streaming results
	resultChan := make(chan benchmark.BenchmarkResult, 100)
	doneChan := make(chan struct{})

	// Start benchmark in goroutine
	go func() {
		defer close(doneChan)
		runnerMutex.Lock()
		runner := benchmarkRunner
		runnerMutex.Unlock()

		config := benchmark.BenchmarkConfig{
			DataSize:   req.DataSize,
			Structures: req.Structures,
			Operation:  req.Operation,
		}

		runner.RunBenchmark(config, func(result benchmark.BenchmarkResult) {
			select {
			case resultChan <- result:
			default:
				// Channel full, skip this update
			}
		})
	}()

	// Stream results
	clientGone := c.Request.Context().Done()
	completedCount := 0
	totalStructures := len(req.Structures)

	for {
		select {
		case <-clientGone:
			benchmarkRunner.Stop()
			return
		case result := <-resultChan:
			data, _ := json.Marshal(result)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()

			if result.Completed {
				completedCount++
				if completedCount >= totalStructures {
					fmt.Fprintf(c.Writer, "event: complete\ndata: {\"message\": \"All benchmarks completed\"}\n\n")
					c.Writer.Flush()
					return
				}
			}
		case <-doneChan:
			fmt.Fprintf(c.Writer, "event: complete\ndata: {\"message\": \"All benchmarks completed\"}\n\n")
			c.Writer.Flush()
			return
		}
	}
}

// HandleStopBenchmark stops any running benchmark
func HandleStopBenchmark(c *gin.Context) {
	runnerMutex.Lock()
	benchmarkRunner.Stop()
	benchmarkRunner = benchmark.NewRunner()
	runnerMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Benchmark stopped",
	})
}

// HandleBenchmarkStatus returns current benchmark status
func HandleBenchmarkStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"structures": []string{"hashmap", "btree", "rbtree", "avltree"},
		"operations": []string{"insert", "search"},
	})
}
