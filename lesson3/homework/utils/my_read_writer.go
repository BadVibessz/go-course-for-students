package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

type StdinReader interface {
	ReadFromStdin(offset int, limit int) ([]byte, error)
}

type FileReader interface {
	ReadFromFile(filename string, offset int) ([]byte, error)
}

type StdoutWriter interface {
	WriteToStdout(output []byte)
}

type FileWriter interface {
	WriteToFile(filename string, bytes []byte) error
}

type StdFileReadWriter interface {
	StdinReader
	FileReader
	StdoutWriter
	FileWriter
}

type MyReadWriter struct {
}

func checkOffset(offset int, arrayLength int) error {

	if offset != 0 && offset >= arrayLength {
		return errors.New("offset is greater than input length")
	}

	if offset < 0 {
		return errors.New("offset is a negative value")
	}
	return nil
}

func (r *MyReadWriter) ReadFromStdin(offset int, limit int) ([]byte, error) {

	// todo: why is not working?
	//var input []byte
	//
	//reader := bufio.NewReader(os.Stdin)
	//_, err := reader.Read(input)

	if limit < -1 {
		return nil, errors.New("limit is a negative value")
	}

	// todo: bad?
	if limit == -1 {
		limit = math.MaxInt64
	}

	// var readSize int
	var input []byte
	in := bufio.NewReader(os.Stdin)
	for {

		if limit == 0 {
			break
		}

		// todo: block size?
		c, err := in.ReadByte()
		if err == io.EOF {

			if offset > 0 {
				return nil, errors.New("offset is greater than input length")
			}
			break
		}

		if offset <= 0 {
			input = append(input, c)
			limit--
		}

		offset--
	}

	return input, nil
}

func (r *MyReadWriter) ReadFromFile(filename string, offset int, limit int) ([]byte, error) {

	// todo: read byte by byte (or by block size)
	input, err := os.ReadFile(filename)

	if err != nil {
		return nil, err // todo: custom error?
	}

	err = checkOffset(offset, len(input))
	if err != nil {
		return nil, err
	}

	if limit == -1 || limit > len(input) {
		limit = len(input) - offset
	}

	if limit < -1 {
		return nil, errors.New("limit is a negative value")
	}

	return input[offset : offset+limit], nil
}

func (r *MyReadWriter) WriteToStdout(output []byte) {
	fmt.Print(bytes.NewBuffer(output).String())
}

func (r *MyReadWriter) WriteToFile(filename string, bytes []byte) error {

	file, err := os.Create(filename)

	if err != nil {
		return err // todo: custom error?
	}

	defer file.Close()

	_, err2 := file.Write(bytes)

	if err2 != nil {
		return err2 // todo: custom error?
	}
	return nil
}
