package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"lecture03_homework/utils"
	"log"
	"os"
	"slices"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	Conv      string
	Uppercase bool
	Lowercase bool
	Trim      bool
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.StringVar(&opts.Conv, "conv", "", "conversion rule. by default - none")

	flag.IntVar(&opts.Offset, "offset", 0, "offset of the input. by default - 0")
	flag.IntVar(&opts.Limit, "limit", -1, "limit of the input size. by default - -1")

	flag.Parse()

	return &opts, nil
}

func validateFlags(opts *Options) error {

	if opts.From != "" { // cannot read from non-existing file

		if _, err := os.Stat(opts.From); err != nil && errors.Is(err, os.ErrNotExist) {
			// file does not exist
			return errors.New("can not read from non-existing file")
		}
	}

	if opts.To != "" { // cannot write to existing file

		if _, err := os.Stat(opts.To); err == nil {
			// file exists
			return errors.New("can not write to existing file")
		}
	}

	if opts.Limit < -1 {
		return errors.New("limit can not be negative value")
	}

	if opts.Offset <= -1 {
		return errors.New("offset can not be negative value")
	}

	// convs

	convs := strings.Split(opts.Conv, ",")
	for _, conv := range convs {

		conv = strings.TrimSpace(conv)

	}

	upInd := slices.IndexFunc(convs, func(s string) bool { return strings.ToLower(strings.TrimSpace(s)) == "upper_case" })
	lowInd := slices.IndexFunc(convs, func(s string) bool { return strings.ToLower(strings.TrimSpace(s)) == "lower_case" })
	trimInd := slices.IndexFunc(convs, func(s string) bool { return strings.ToLower(strings.TrimSpace(s)) == "trim_spaces" })

	if upInd != -1 && lowInd != -1 {
		return errors.New("uppercase and lowercase can not live together")
	}

	if upInd != -1 {
		opts.Uppercase = true
	}

	if lowInd != -1 {
		opts.Lowercase = true
	}

	if trimInd != -1 {
		opts.Trim = true
	}

	if upInd == -1 && lowInd == -1 && trimInd == -1 && len(convs) != 0 {
		return errors.New("wrong conv flag provided")
	}

	return nil
}

func main() {
	opts, err := ParseFlags()

	logger := log.New(os.Stderr, "", 3)

	var rw = utils.MyReadWriter{}
	var tc = utils.MyTextConverter{}

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	err = validateFlags(opts)
	if err != nil {
		logger.Println("error occurred during validating flags:", err)
		os.Exit(1)
	}

	var input []byte

	if opts.From == "" {

		input, err = rw.ReadFromStdin(opts.Offset, opts.Limit)
		if err != nil {
			logger.Println("error occurred during reading from stdin:", err)
			os.Exit(1)
		}

	} else {

		input, err = rw.ReadFromFile(opts.From, opts.Offset, opts.Limit)
		if err != nil {
			logger.Println("error occurred during reading from file:", err)
			os.Exit(1)
		}

	}

	if opts.Uppercase {
		str := bytes.NewBuffer(input).String()
		input = []byte(tc.Uppercase(str))
	}

	if opts.Lowercase {
		str := bytes.NewBuffer(input).String()
		input = []byte(tc.Lowercase(str))
	}

	if opts.Trim {
		str := bytes.NewBuffer(input).String()
		input = []byte(tc.TrimSpaces(str))
	}

	if opts.To == "" {
		rw.WriteToStdout(input)
	} else {

		err = rw.WriteToFile(opts.To, input)
		if err != nil {
			logger.Println("error occurred during writing to file:", err)
			os.Exit(1)
		}
	}

}
