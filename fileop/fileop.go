package fileop

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func LineCount(filename string) (cnt int, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(filename); err != nil {
		return
	}
	defer file.Close()

	var (
		buf     []byte
		n       int
		lineSep = []byte{'\n'}
	)
	buf = make([]byte, 32768) //32k
	for {
		n, err = file.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		cnt += bytes.Count(buf[:n], lineSep)

		if err == io.EOF {
			err = nil
			return
		}
	}
	return
}

func YALineCount(filename string) (cnt int, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(filename); err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cnt++
	}

	err = scanner.Err()
	return
}
