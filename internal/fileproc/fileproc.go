// Package fileproc counts lines and bytes of files, isolating a panic in one
// file from the rest of the batch.
package fileproc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Result is the outcome of processing one path.
type Result struct {
	Path     string
	Lines    int
	Bytes    int
	Err      error // non-nil if the file failed
	Panicked bool  // true if a panic was recovered for this file
}

func Process(paths []string) []Result {
	results := make([]Result, 0, len(paths))
	for _, p := range paths {
		results = append(results, processFile(p))
	}
	return results
}

func processFile(path string) (res Result) {
	defer func() {
		if r := recover(); r != nil {
			res = Result{
				Path:     path,
				Panicked: true,
				Err:      fmt.Errorf("fileproc: panic processing %s: %v", path, r),
			}
		}
	}()

	f, err := os.Open(path)
	if err != nil {

		return Result{Path: path, Err: err}
	}
	defer f.Close()

	br := bufio.NewReader(f)

	first, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return Result{Path: path, Err: err}
	}
	if strings.TrimSpace(first) == "PANIC" {
		panic("fileproc: boom")
	}

	size := len(first)
	lines := strings.Count(first, "\n")
	var last byte
	if size > 0 {
		last = first[size-1]
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := br.Read(buf)
		if n > 0 {
			size += n
			lines += bytes.Count(buf[:n], []byte{'\n'})
			last = buf[n-1]
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return Result{Path: path, Err: err}
		}
	}

	if size > 0 && last != '\n' {
		lines++
	}
	return Result{Path: path, Lines: lines, Bytes: size}
}
