# ZFS Benchmark Inefficiency Test

This Go program is designed to test potential inefficiencies in ZFS storage performance. The benchmark focuses on read and write performance with varying block sizes, data sizes, and concurrency levels (goroutines). It also supports customizable sync frequency, allowing you to simulate different real-world scenarios.

## Features
- Customizable block size, data size, and sync frequency
- Concurrent read/write operations using multiple goroutines
- Logs performance results to a file (`benchmark.log`)
- Cleans up test files after benchmarking

## Requirements
- Go (Golang) installed on your system
- ZFS or any storage system to run the benchmark

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/zfs-benchmark.git
   ```
2. Navigate to the directory:
   ```bash
   cd zfs-benchmark
   ```

3. Run `go mod tidy` to clean up any unused dependencies and ensure your Go module is tidy:
   ```bash
   go mod tidy
   ```

## Building the Application

Before running the application, you need to build the Go binary. Run the following command to build the executable:

```bash
go build -o zfs-benchmark src/main.go
```

This will create an executable named `zfs-benchmark` in the current directory.

## Running the Application

Once the application is built, you can run the executable using the following command:

```bash
./zfs-benchmark [flags]
```

### Flags
- `-blocksize`: Size of blocks in bytes. (Default: 4096 bytes)
- `-datasize`: Total size of data to write in bytes. (Default: 128 MiB)
- `-syncfreq`: Number of blocks to write before performing a sync (flush). (Default: 10,000)
- `-goroutines`: Number of concurrent goroutines for read/write operations. (Default: 4)

### Example Usage

1. **Default settings**:
   ```bash
   ./zfs-benchmark
   ```
   This runs the test with:
   - Block size: 4096 bytes
   - Data size: 128 MiB
   - Sync frequency: 10,000 blocks
   - Goroutines: 4

2. **Custom block size, data size, and sync frequency**:
   ```bash
   ./zfs-benchmark -blocksize=8192 -datasize=256000000 -syncfreq=5000
   ```
   This runs the test with:
   - Block size: 8192 bytes
   - Data size: 256 MiB
   - Sync frequency: 5,000 blocks

3. **Increasing concurrency**:
   ```bash
   ./zfs-benchmark -blocksize=1024 -datasize=256000000 -syncfreq=5000 -goroutines=8
   ```
   This runs the test with:
   - Block size: 1024 bytes
   - Data size: 256 MiB
   - Sync frequency: 5,000 blocks
   - Goroutines: 8 (8 concurrent read/write threads)

## Clean-Up

The program automatically removes all test files after benchmarking. However, if you want to manually clean up the files, you can remove them as follows:

```bash
rm testfile_* 
```

## License

This project is licensed under the BSD 3-Clause License.

---

## Benchmark Tests and Results

This section summarizes the results from running the benchmark with different configurations.

### Test 1: 4 KiB Block Size (Default)
```bash
./zfs-benchmark
```
- **Write Performance**: 4.86s
- **Read Performance**: 147.42ms

### Test 2: 8 KiB Block Size
```bash
./zfs-benchmark -blocksize=8192 -datasize=256000000 -syncfreq=5000
```
- **Write Performance**: 19.75s
- **Read Performance**: 377.55ms

### Test 3: 1 KiB Block Size
```bash
./zfs-benchmark -blocksize=1024 -datasize=256000000 -syncfreq=5000
```
- **Write Performance**: 18.18s
- **Read Performance**: 1.60s

### Test 4: 2 KiB Block Size
```bash
./zfs-benchmark -blocksize=2048 -datasize=256000000 -syncfreq=5000
```
- **Write Performance**: 13.10s
- **Read Performance**: 1.42s

### Summary of Results

| Block Size | Write Time | Read Time    |
|------------|------------|--------------|
| 4 KiB      | 4.86s      | 147.42ms     |
| 8 KiB      | 19.75s     | 377.55ms     |
| 1 KiB      | 18.18s     | 1.60s        |
| 2 KiB      | 13.10s     | 1.42s        |

### Conclusion
The **4 KiB block size** provides the best overall performance for both reads and writes. It is optimal for ZFS performance under the tested conditions.
