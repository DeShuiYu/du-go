package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func calcAllFileSize(root string) (string, int64) {
	var fileSize int64 = 0
	var mutex sync.Mutex
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() {
			info, _ := dir.Info()
			mutex.Lock()
			fileSize += info.Size()
			mutex.Unlock()
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return root, fileSize
}

func FormatSize(bytes int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func main() {

	//basePath := os.Args[1]
	basePath := "/Users/dsy/workplaces/"
	entries, _ := os.ReadDir(basePath)
	var wg sync.WaitGroup
	for _, entry := range entries {
		wg.Add(1)
		go func(e os.DirEntry) {
			defer wg.Done()
			if e.IsDir() {
				root, fileSize := calcAllFileSize(filepath.Join(basePath, e.Name()))
				fmt.Printf("%s,%s\n", root, FormatSize(fileSize))
			} else {
				info, _ := e.Info()
				fmt.Printf("%s,%s\n", filepath.Join(basePath, e.Name()), FormatSize(info.Size()))
			}
		}(entry)
	}
	wg.Wait()
}
