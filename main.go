package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

func CalcAllFileSizeSum(root string) (string, int64) {
	var fileSize int64 = 0
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() {
			info, _ := dir.Info()

			fileSize += info.Size()
		}
		return nil
	})
	if err != nil {
		fmt.Errorf("%s error: %v", root, err)
	}
	return root, fileSize
}

func main() {
	st := time.Now()
	basePath := os.Args[1]
	entries, _ := os.ReadDir(basePath)
	exclude := []string{"share"}
	var wg sync.WaitGroup
	for _, entry := range entries {
		if slices.Contains(exclude, entry.Name()) {
			fmt.Println("Skipping", entry.Name())
			continue
		}

		wg.Add(1)
		go func(e os.DirEntry) {
			defer wg.Done()
			fullpath := filepath.Join(basePath, e.Name())
			if e.IsDir() {
				_, fileSize := CalcAllFileSizeSum(fullpath)
				fmt.Printf("%s,%s\n", fullpath, humanize.Bytes(uint64(fileSize)))
			} else {
				info, _ := e.Info()
				fmt.Printf("%s,%s\n", fullpath, humanize.Bytes(uint64(info.Size())))
			}
		}(entry)
	}
	wg.Wait()
	fmt.Println("Total time:", time.Since(st))
}
