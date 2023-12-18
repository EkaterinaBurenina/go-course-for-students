package storage

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func walkDir(g *errgroup.Group, ctx context.Context, dir Dir, fSizeList chan<- int64) error {
	childDirs, files, err := dir.Ls(ctx)
	if err != nil {
		fmt.Println("dir ls ERROR: ", err)
		return err
	}

	for _, childDir := range childDirs {
		childDir := childDir
		g.Go(func() error {
			err := walkDir(g, ctx, childDir, fSizeList)
			return err
		})
	}
	for _, f := range files {
		fSize, err := f.Stat(ctx)
		if err != nil {
			//fmt.Println("f stat ERROR: ", err)
			return err
		}
		fSizeList <- fSize
	}
	return nil
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	g, ctx := errgroup.WithContext(ctx)
	fSizeList := make(chan int64)
	res := Result{}
	done := make(chan struct{})
	go func() {
		for size := range fSizeList {
			res.Size += size
			res.Count++
		}
		close(done)
	}()

	g.Go(func() error {
		err := walkDir(g, ctx, d, fSizeList)
		return err
	})

	if err := g.Wait(); err != nil {
		return res, err
	}

	close(fSizeList)
	<-done

	return res, nil
}
