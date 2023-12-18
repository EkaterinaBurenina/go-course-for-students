package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

//func GetFunctionName(i interface{}) string {
//	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
//}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {

	streamHandler := func(ctx context.Context, valueStreamIn In, i int) In {
		//fmt.Println("stage", i, GetFunctionName(stages[i]), "start")
		valueStreamOut := make(chan any)
		go func() {
			defer close(valueStreamOut)
			for {
				select {
				case <-ctx.Done():
					//fmt.Println("ctx closed")
					return
				case v, ok := <-valueStreamIn:
					if !ok {
						return
					} else {
						//fmt.Println("stage", i, GetFunctionName(stages[i]), "in value:", v)
						valueStreamOut <- v
					}
				}
			}
		}()
		return valueStreamOut
	}

	for i, stage := range stages {
		in = stage(streamHandler(ctx, in, i))
	}

	return in
}
