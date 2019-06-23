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
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("couldn't open file '%s': %v", file, err)
	}
	defer f.Close()

	lines := uniques(f)

	f, err = os.Create(file)
	if err != nil {
		return fmt.Errorf("couldn't create file '%s': %v", file, err)
	}
	defer f.Close()

	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("couldn't close file '%s': %v", file, err)
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
