package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From string
	To   string

	Offset    int
	Limit     int
	BlockSize int

	Conv string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")

	flag.IntVar(&opts.Offset, "offset", 0, "offset in bytes")
	flag.IntVar(&opts.Limit, "limit", 0, "limit in bytes")
	flag.IntVar(&opts.BlockSize, "block-size", 1024000000, "block size in bytes")

	flag.StringVar(&opts.Conv, "conv", "", "convert to format")

	flag.Parse()

	return &opts, nil
}

type Reader interface {
	Read(p []byte) ([]byte, error)
}

type MyReader struct {
	r    io.Reader
	conv []string
}

// type FileReader struct {
// 	r    io.Reader
// 	conv []string
// }

func NewReader(s *os.File, offset int, limit int) io.Reader {
	r := bufio.NewReader(s)
	size := r.Size()

	if offset > size {
		fmt.Fprintln(os.Stderr, "offset is bigger than file size")
		os.Exit(1)
	}
	_, err := r.Discard(int(offset))
	if err != nil {
		fmt.Fprintln(os.Stderr, "can not discard the following offset bytes:", err)
		os.Exit(1)
	}
	if limit != 0 {
		return io.LimitReader(r, int64(limit))
	}
	return r
}

func (r MyReader) Read(p []byte) (res []byte, err error) {
	fmt.Println("\nsource len: ", len(p))
	n, err := r.r.Read(p)
	fmt.Println("\n n: ", n)
	res = p[:n]
	// fmt.Println("\n res input len: ", len(res))

	have_upper_case := false
	have_lower_case := false

	for _, v := range r.conv {
		switch v {
		case "upper_case":
			if have_lower_case {
				fmt.Fprintln(os.Stderr, "can not convert to lower case after upper case")
				os.Exit(1)
			}
			have_upper_case = true
			res = bytes.ToUpper(res)
		case "lower_case":
			if have_upper_case {
				fmt.Fprintln(os.Stderr, "can not convert to lower case after upper case")
				os.Exit(1)
			}
			have_lower_case = true
			res = bytes.ToLower(res)
		case "trim_spaces":
			res = bytes.TrimSpace(res)
			// for _, v := range res {
			// 	fmt.Println("v: ", string(v), ", is_space: ", unicode.IsSpace(rune(v)))
			// }
		}
	}

	return
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	conv := strings.Split(opts.Conv, ",")

	var r Reader

	if opts.From == "" {
		r = MyReader{r: NewReader(os.Stdin, opts.Offset, opts.Limit), conv: conv}
	} else {
		file, err := os.Open(opts.From)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can not open file:", err)
			os.Exit(1)
		}
		r = MyReader{r: NewReader(file, opts.Offset, opts.Limit), conv: conv}
		defer file.Close()
	}

	var w *bufio.Writer
	if opts.To == "" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		file, err := os.OpenFile(opts.To, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can not open file:", err)
			os.Exit(1)
		}
		w = bufio.NewWriter(file)
		defer file.Close()
	}

	buf := make([]byte, opts.BlockSize)
	fmt.Println("\nbuf size: ", len(buf))
	for {
		res, err := r.Read(buf)
		if err == io.EOF {
			break
		}

		_, err = w.Write(res)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can not write to file:", err)
			os.Exit(1)
		}
		err = w.Flush()
		if err != nil {
			fmt.Fprintln(os.Stderr, "can not flush to file:", err)
			os.Exit(1)
		}
	}
}
