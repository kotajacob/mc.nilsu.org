package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// keySlice is a slice of entries which contain key and value strings. They may
// be parsed from a file with one entry per line with the format `key;value\n`.
type keySlice []entry

// entry holds key and value strings as part of a keySlice.
type entry struct {
	Key   string
	Value string
}

// parseKeyFile reads a file of keys and values in the format `key:value\n` and
// stores it in a keySlice. If a line does not have 2 fields or otherwise cannot
// be read an error is returned.
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
		if len(t) != 2 {
			return nil, fmt.Errorf("failed parsing %v: "+
				"invalid number of fields: %v\n", name, scanner.Text())
		}

		e := entry{
			Key:   t[0],
			Value: t[1],
		}
		entries = append(entries, e)
	}
	return entries, scanner.Err()
}
