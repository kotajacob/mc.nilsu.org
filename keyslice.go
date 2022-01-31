package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// keySlice is a slice of entries which contain key, value, and optional url
// strings. They may be parsed from a file with one entry per line with the
// format `key;value;url\n`.
type keySlice []entry

// entry holds key value and url strings as part of a keySlice.
type entry struct {
	Key   string
	Value string
	URL   string
}

// parseKeyFile reads a file of keys, values, and urls in the format
// `key:value;url\n` and stores it in a keySlice. If a line has less than 2
// fields or otherwise cannot be read an error is returned.
func parseKeyFile(name string) (keySlice, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries keySlice

	for scanner.Scan() {
		t := strings.Split(scanner.Text(), ";")
		var e entry

		if len(t) >= 2 {
			e.Key = t[0]
			e.Value = t[1]
		} else {
			return nil, fmt.Errorf("failed parsing %v: "+
				"invalid number of fields: %v\n", name, scanner.Text())
		}
		if len(t) == 3 {
			e.URL = t[2]
		}

		entries = append(entries, e)
	}
	return entries, scanner.Err()
}
