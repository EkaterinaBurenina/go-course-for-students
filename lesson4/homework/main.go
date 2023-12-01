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

func walkDir2(wg *sync.WaitGroup, path string, f_size_list chan int) {
	// defer wg.Done()
	fmt.Println("WALK ", path)

	items, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fmt.Println(item.Name(), "is_dir: ", item.IsDir())

		if item.IsDir() {
			wg.Add(1)

			go walkDir2(wg, path+"/"+item.Name(), f_size_list)
			defer wg.Done()
			// go func() {
			// 	defer wg.Done()
			// 	walkDir2(wg, path+"/"+item.Name(), f_size_list)
			// }()
		} else {
			f_info, _ := item.Info()
			f_size_list <- int(f_info.Size())
		}
	}
}

func main() {
	path := "."
	fmt.Println("sync result: ", walkDir(path))

	var result int

	wg := sync.WaitGroup{}
	f_size_list := make(chan int)
	defer close(f_size_list)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for size := range f_size_list {
			result += size
			fmt.Println("result += size: ", result)
		}
	}()

	wg.Add(1)
	// go walkDir2(&wg, path, f_size_list)
	go func() {
		defer wg.Done()
		walkDir2(&wg, path, f_size_list)
	}()

	wg.Wait()

	fmt.Println("async result: ", result)
}
