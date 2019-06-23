package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func uniques(r io.Reader) []string {
	lines := make([]string, 0, 10)
	seen := make(map[string]struct{})

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if _, ok := seen[line]; !ok {
			seen[line] = struct{}{}
			lines = append(lines, line)
		}
	}
	return lines
}

func uniqueFile(file string) error {
	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %v", file, err)
	}
	defer f.Close()

	lines := uniques(f)
	// write unique lines
	if _, err = f.Seek(io.SeekStart, 0); err != nil {
		return fmt.Errorf("couldn't jump to beginning of %s: %v", file, err)
	}
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	// set new filesize
	pos, err := f.Seek(io.SeekStart, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("couldn't determine position of %s: %v", file, err)
	}
	if err = f.Truncate(pos); err != nil {
		return fmt.Errorf("couldn't truncate file %s: %v", file, err)
	}
	return nil
}

func main() {
	if len(os.Args) > 1 {
		// read and write files in succession
		for _, f := range os.Args[1:] {
			if err := uniqueFile(f); err != nil {
				log.Fatalf("error while handling '%s': %v", f, err)
			}
		}
	} else {
		// read stdin, write stdout
		lines := uniques(os.Stdin)
		for _, l := range lines {
			fmt.Println(l)
		}
	}
}
