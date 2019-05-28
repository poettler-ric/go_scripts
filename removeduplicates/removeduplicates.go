package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func uniques(r io.Reader) []string {
	lines := make([]string, 0, 10)
	sortedLines := make([]string, 0, 10)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		i := sort.SearchStrings(sortedLines, line)
		if i == len(sortedLines) || sortedLines[i] != line {
			// line not found - add it
			sortedLines = append(sortedLines, "")
			copy(sortedLines[i+1:], sortedLines[i:])
			sortedLines[i] = line

			lines = append(lines, line)
		}
	}
	return lines
}

func uniqueFile(file string) {
	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("couldn't open file %s: %v", file, err)
	}
	defer f.Close()

	lines := uniques(f)
	// write unique lines
	if _, err = f.Seek(io.SeekStart, 0); err != nil {
		log.Fatalf("couldn't jump to beginning of %s: %v", file, err)
	}
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	// set new filesize
	pos, err := f.Seek(io.SeekStart, io.SeekCurrent)
	if err != nil {
		log.Fatalf("couldn't determine position of %s: %v", file, err)
	}
	if err = f.Truncate(pos); err != nil {
		log.Fatalf("couldn't truncate file %s: %v", file, err)
	}
}

func main() {
	if len(os.Args) > 1 {
		// read and write files in succession
		for _, f := range os.Args[1:] {
			uniqueFile(f)
		}
	} else {
		// read stdin, write stdout
		lines := uniques(os.Stdin)
		for _, l := range lines {
			fmt.Println(l)
		}
	}
}