package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

var (
	root     string
	execlude string
)

func GetDirOrFileDiskUsage(root string) int64 {
	totalSize := int64(0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, _ := d.Info()
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		fmt.Errorf("%v", err)
	}
	return totalSize
}

func main() {
	st := time.Now()
	flag.StringVar(&root, "p", "", "输入需要扫描的文件夹")
	flag.StringVar(&execlude, "e", "", "输入需要排除的文件或者文件夹")
	flag.Parse()
	execlude_array := strings.Split(execlude, ",")
	entries, _ := os.ReadDir(root)

	if root == "" {
		panic("请输入需要扫描的文件夹")
	}
	fmt.Printf("当前需要排除的文件夹或者文件为:%v\n", execlude_array)

	var wg sync.WaitGroup
	for _, entry := range entries {
		filecurrentpath := filepath.Join(root, entry.Name())
		if slices.Contains(execlude_array, entry.Name()) || slices.Contains(execlude_array, filecurrentpath) {
			fmt.Printf("%s跳过<%s>\n", "\033[31m", filecurrentpath)
			continue
		}
		wg.Add(1)
		go func(e os.DirEntry, currentpath string) {
			defer wg.Done()
			filesize := GetDirOrFileDiskUsage(currentpath)
			if filesize >= humanize.GByte {
				fmt.Printf("%s%s,%s\n", "\033[32m", currentpath, humanize.Bytes(uint64(filesize)))
			}
		}(entry, filecurrentpath)
	}
	wg.Wait()
	fmt.Printf("%stotal time:%v", "\033[0m", time.Since(st))
}
