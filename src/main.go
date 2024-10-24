package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	DefaultTestFileName = "testfile.bin"
	DefaultBlockSize    = 4096  // Default block size in bytes (4 KiB)
	DefaultDataSize     = 128 * 1024 * 1024 // Default data size of 128 MiB
	DefaultSyncFreq     = 10000 // Default sync frequency
	DefaultNumGoroutines = 4    // Default number of concurrent goroutines
)

// Struct for configuration
type Config struct {
	blockSize    int
	totalDataSize int
	syncFrequency int
	numGoroutines int
}

func main() {
	// Command-line flags
	blockSize := flag.Int("blocksize", DefaultBlockSize, "Size of blocks in bytes")
	totalDataSize := flag.Int("datasize", DefaultDataSize, "Total size of data to write (in bytes)")
	syncFrequency := flag.Int("syncfreq", DefaultSyncFreq, "Number of blocks before flushing (sync)")
	numGoroutines := flag.Int("goroutines", DefaultNumGoroutines, "Number of concurrent goroutines")
	flag.Parse()

	config := Config{
		blockSize:    *blockSize,
		totalDataSize: *totalDataSize,
		syncFrequency: *syncFrequency,
		numGoroutines: *numGoroutines,
	}

	fmt.Println("ZFS Benchmark Inefficiency Test with concurrency")

	// Open log file
	logFile, err := os.OpenFile("benchmark.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	// Concurrent Write Test
	concurrentWriteTest(config, logFile)

	// Concurrent Read Test
	concurrentReadTest(config, logFile)

	// Clean up the test file
	cleanUpTestFile()
}

// concurrentWriteTest performs concurrent writes with multiple goroutines
func concurrentWriteTest(config Config, logFile *os.File) {
	fmt.Printf("\nStarting Concurrent Write Performance Test with Block Size: %d bytes, Goroutines: %d\n", config.blockSize, config.numGoroutines)

	var wg sync.WaitGroup
	start := time.Now()

	// Launch goroutines to write data concurrently
	for i := 0; i < config.numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			filename := fmt.Sprintf("%s_%d", DefaultTestFileName, id)
			file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", filename, err)
				return
			}
			defer file.Close()

			// Generate dummy data to write
			data := make([]byte, config.blockSize)

			// Write data in blocks with sync every `syncFrequency`
			for j := 0; j < config.totalDataSize/config.blockSize; j++ {
				_, err := file.Write(data)
				if err != nil {
					fmt.Printf("Error writing to file %s: %v\n", filename, err)
					return
				}

				// Flush to disk every syncFrequency blocks
				if j%config.syncFrequency == 0 {
					err = file.Sync()
					if err != nil {
						fmt.Printf("Error syncing to disk %s: %v\n", filename, err)
						return
					}
				}
			}

			// Ensure all data is written to disk
			err = file.Sync()
			if err != nil {
				fmt.Printf("Error syncing final data to disk %s: %v\n", filename, err)
				return
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Concurrent Write Performance Test completed in: %v\n", duration)
	logResults(logFile, "Write", config.blockSize, config.totalDataSize, config.numGoroutines, duration)
}

// concurrentReadTest performs concurrent reads with multiple goroutines
func concurrentReadTest(config Config, logFile *os.File) {
	fmt.Printf("\nStarting Concurrent Read Performance Test with Block Size: %d bytes, Goroutines: %d\n", config.blockSize, config.numGoroutines)

	var wg sync.WaitGroup
	start := time.Now()

	// Launch goroutines to read data concurrently
	for i := 0; i < config.numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			filename := fmt.Sprintf("%s_%d", DefaultTestFileName, id)
			file, err := os.Open(filename)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", filename, err)
				return
			}
			defer file.Close()

			// Create a buffer for reading in chunks
			buf := make([]byte, config.blockSize)

			// Read data in blocks
			for {
				n, err := file.Read(buf)
				if n == 0 || err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("Error reading file %s: %v\n", filename, err)
					return
				}
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Concurrent Read Performance Test completed in: %v\n", duration)
	logResults(logFile, "Read", config.blockSize, config.totalDataSize, config.numGoroutines, duration)
}

// logResults writes the test results to the log file
func logResults(logFile *os.File, testType string, blockSize, dataSize, numGoroutines int, duration time.Duration) {
	logLine := fmt.Sprintf("%s Test - Block Size: %d, Data Size: %d, Goroutines: %d, Duration: %v\n",
		testType, blockSize, dataSize, numGoroutines, duration)
	_, err := logFile.WriteString(logLine)
	if err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

// cleanUpTestFile removes the test files after benchmarking
func cleanUpTestFile() {
	for i := 0; i < DefaultNumGoroutines; i++ {
		filename := fmt.Sprintf("%s_%d", DefaultTestFileName, i)
		err := os.Remove(filename)
		if err != nil {
			fmt.Printf("Error cleaning up test file %s: %v\n", filename, err)
			return
		}
	}
	fmt.Println("Test files removed successfully.")
}

