package main

import (
	"flag"
	"fmt"
	"lecture03_homework/utils"
	"log"
	"os"
)

type Options struct {
	From   string
	To     string
	Offset int
	Limit  int
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "offset of the input. by default - 0")
	flag.IntVar(&opts.Limit, "limit", -1, "limit of the input size. by default - -1")

	flag.Parse()

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	logger := log.New(os.Stderr, "", 3)
	var rw = utils.MyReadWriter{}

	var input []byte

	if opts.From == "" {

		input, err = rw.ReadFromStdin(opts.Offset, opts.Limit)
		if err != nil {
			// todo: handle
			logger.Println("Error occurred during reading from stdin:", err)
		}

	} else {

		input, err = rw.ReadFromFile(opts.From, opts.Offset, opts.Limit)
		if err != nil {
			// todo: handle
			logger.Println("Error occurred during reading from file:", err)

		}

	}

	if opts.To == "" {
		rw.WriteToStdout(input)
	} else {

		err = rw.WriteToFile(opts.To, input)
		if err != nil {
			// todo: handle
			logger.Println("Error occurred during writing to file:", err)
		}
	}
}
