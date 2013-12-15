package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"io"
	"sort"
	"path/filepath"
	"strings"
)

func main() {
	// parse command line flags

	var (
		listDir string
		listName string
		removeItem int
	)

	flag.StringVar(&listDir, "d", "~/Documents", "directory of list")
	flag.StringVar(&listName, "n", "to.txt", "filename of list")
	flag.IntVar(&removeItem, "r", -1, "item to be removed")
	flag.Parse()

	// sanitize path to list

	tildeAt := strings.Index(listDir, "~")
	if tildeAt != -1 {
		homeDir := os.Getenv("HOME")
		if len(homeDir) == 0 {
			panic("can't determine home directory")
		}
		listDir = strings.Replace(listDir, "~", homeDir, 1)
	}

	listDir, err := filepath.Abs(listDir)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(listDir); os.IsNotExist(err) {
		panic(err)
	}

	listPath := filepath.Join(listDir, listName)

	// get items from list

	lines, err := ReadLines(listPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch {
	case removeItem >= 0:
		// remove item at removeItem
		if removeItem >= len(lines) {
			panic("item to remove doesn't exist")
		}
		lines[removeItem] = lines[len(lines) - 1]
		lines = lines[:len(lines) - 1]
		sort.Strings(lines)
		err = WriteLines(listPath, lines)
		if err != nil {
			panic(err)
		}
	case flag.NArg() == 0:
		// display item list
		width := 1
		if len(lines) > 10 {
			width = 2
		}
		var format = fmt.Sprintf("%%%dd - %%s\n", width)
		for i, line := range lines{
			fmt.Printf(format, i, line)
		}
	default:
		// join rest of arguments as item to add
		to := strings.Join(flag.Args(), " ")
		lines = append(lines, to)
		sort.Strings(lines)
		if err := WriteLines(listPath, lines); err != nil {
			panic(err)
		}
	}

	return
}

func ReadLines(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read lines from buffered reader
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\n")
		if len(line) != 0 {
			lines = append(lines, line)
		}
	}

	// suppress EOF error
	if err == io.EOF {
		err = nil
	}

	return
}

func WriteLines(path string, lines []string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	// write each line to file
	for _, line := range lines {
		_, err = f.WriteString(line + "\n");
		if err != nil {
			return
		}
	}
	return
}

