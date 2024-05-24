package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type onChunk = func(chunk []byte, breakFunc func())

func chunkedReader(reader io.Reader, onChunk onChunk) {
	chunk := make([]byte, 8192)
	eof := false
	for !eof {
		r, err := reader.Read(chunk)
		if err != nil && err != io.EOF {
			return
		} else if err == io.EOF || r == 0 {
			eof = true
		}

		onChunk(chunk[:r], func() {
			eof = true
		})
	}
}

func chunkedFile(file *os.File, onChunk onChunk) {
	chunk := make([]byte, 8192)
	eof := false
	for !eof {
		r, err := file.Read(chunk)
		if err != nil && err != io.EOF {
			return
		} else if err == io.EOF {
			eof = true
		}

		onChunk(chunk[:r], func() {
			eof = true
		})
	}
}

func fromString(value string, receiver any) error {
	err := json.Unmarshal([]byte(value), receiver)
	if err != nil {
		return err
	}
	return nil
}

func getPathForJPEG(path string) string {
	absolutePath, _ := os.Getwd()
	return fmt.Sprintf("%s/files/%s.jpeg", absolutePath, path)
}

func getPathForName(path string) string {
	absolutePath, _ := os.Getwd()
	return fmt.Sprintf("%s/files/%s", absolutePath, path)
}

func stringToInt64(str string) (int64, error) {
	n, e := strconv.Atoi(str)
	return int64(n), e
}
