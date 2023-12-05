package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func walkDir(path string) (size int) {
	// fmt.Println("walk ", path)

	items, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		// fmt.Println(item.Name(), "is_dir: ", item.IsDir())

		if item.IsDir() {
			size += walkDir(path + "/" + item.Name())
			// fmt.Println("size += dir_size: ", size)
		} else {
			f_info, _ := item.Info()
			size += int(f_info.Size())
			// fmt.Println("size += file_size: ", size)
		}

	}

	return
}

func walkDir2(wg *sync.WaitGroup, path string, fSizeList chan int) {
	fmt.Println("WALK ", path)

	items, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fmt.Println(item.Name(), "is_dir: ", item.IsDir())

		if item.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				walkDir2(wg, path+"/"+item.Name(), fSizeList)
			}()
		} else {
			f_info, _ := item.Info()
			//fmt.Println(f_info.Size())
			fSizeList <- int(f_info.Size())
		}
	}
}

func main() {
	path := "."
	fmt.Println("sync result: ", walkDir(path))

	var result int

	wg := sync.WaitGroup{}
	fSizeList := make(chan int)

	done := make(chan struct{})
	go func() {
		for size := range fSizeList {
			result += size
			//fmt.Println("result += size: ", result)
		}
		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		walkDir2(&wg, path, fSizeList)
	}()

	wg.Wait()

	close(fSizeList)

	<-done
	fmt.Println("async result: ", result)
}
