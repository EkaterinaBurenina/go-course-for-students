package main

import (
	"context"
	"fmt"
	"time"

	"homework/executor"
)

func PlusOne(in executor.In) executor.Out {
	//fmt.Println("PlusOne start")
	out := make(chan any, len(in))
	go func() {
		//fmt.Println("PlusOne input go routine start")
		defer close(out)
		for item := range in {
			//fmt.Println("PlusOne source:", item, ", result:", item.(int)+1)
			time.Sleep(time.Second * 3)
			out <- item.(int) + 1
		}
	}()
	return out
}

func MultiplyByTwo(in executor.In) executor.Out {
	//fmt.Println("MultiplyByTwo start")
	out := make(chan any, len(in))
	go func() {
		//fmt.Println("MultiplyByTwo input go routine start")
		defer close(out)
		for item := range in {
			//fmt.Println("MultiplyByTwo source:", item, ", result:", item.(int)*2)
			time.Sleep(time.Second * 1)
			out <- item.(int) * 2
		}
	}()
	return out
}

func main() {
	in_cnt := 5
	in := make(chan any, in_cnt)
	go func() {
		for i := 0; i < in_cnt; i++ {
			in <- i
		}
		close(in)
	}()

	result := make([]any, 0, len(in))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*16)
	defer cancel()

	start := time.Now()
	for r := range executor.ExecutePipeline(ctx, in, MultiplyByTwo,
		PlusOne,
		MultiplyByTwo,
		//MultiplyByTwo,
		//PlusOne,
	) {
		fmt.Println("main res:", r)
		result = append(result, r)
	}
	elapsed := time.Since(start)

	fmt.Println("main result:", result, "elapsed:", elapsed)
}
