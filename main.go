package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"io"
	"sort"
	"path/filepath"
	"strings"
)

var dir_default = "/Users/matt/Documents"
var dir_desc = "directory to store list"

var name_default = "to.txt"
var name_desc = "name of list"

var remove_desc = "prefix of item to remove"

var to_dir = flag.String("d", dir_default, dir_desc)
var to_name = flag.String("n", name_default, name_desc)
var should_remove = flag.String("r", "", name_desc)
// print path option

func main() {
	flag.Parse()

	to_path := filepath.Join(*to_dir, *to_name)

	var (
		lines []string
		err error
	)
	if lines, err = read_lines(to_path); err != nil {
		fmt.Print(err)
		return
	}

	if len(*should_remove) != 0 {
		// find item with prefix
		to_remove := sort.Search(len(lines), func(i int) bool { return lines[i] >= *should_remove || strings.HasPrefix(lines[i], *should_remove) })
		// prefix wasn't found or "found item" doesn't have prefix
		if to_remove == len(lines) || !(strings.HasPrefix(lines[to_remove], *should_remove)) {
			fmt.Println("ERROR: no item with that prefix")
			os.Exit(1)
		// multiple items have prefix
		} else if to_remove + 1 < len(lines) && strings.HasPrefix(lines[to_remove + 1], *should_remove) {
			fmt.Println("ERROR: more than one item with that prefix")
			os.Exit(1)
		}

		// remove element at to_remove
		copy(lines[to_remove:], lines[(to_remove + 1):len(lines)])
		lines = lines[:len(lines) - 1]
		if err = write_lines(to_path, lines); err != nil {
			fmt.Print(err)
			return
		}
		return
	}

	// just print the lines
	if len(flag.Args()) == 0 {
		for _, line := range lines {
			fmt.Println("-", line)
		}
		return
	}

	// join rest of arguments as item to add
	to := strings.Join(flag.Args(), " ")
	lines = append(lines, to)
	sort.Strings(lines)
	if err = write_lines(to_path, lines); err != nil {
		fmt.Print(err)
		return
	}
	return
}

// returns all lines in file at path (even blank ones)
func read_lines(path string) (lines []string, err error) {
	var (
		file *os.File
		line string
	)
	// open file for reading
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	// start reading from file
	reader := bufio.NewReader(file)
	for {
		if line, err = reader.ReadString('\n'); err != nil {
			break
		}
		// don't include newline
		line = strings.TrimRight(line, "\n")
		// skip blank lines
		if len(line) != 0 {
			lines = append(lines, strings.TrimRight(line, "\n"))
		}
	}
	// don't report an end of file error
	if err == io.EOF {
		err = nil
	}
	return
}

func write_lines(path string, lines []string) (err error) {
	var (
		file *os.File
	)
	// open file for writing
	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	// write each line to file
	for _, line := range lines {
		_, err = file.WriteString(line + "\n");
		if err != nil {
			return
		}
	}
	return
}

