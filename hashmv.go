// Copyright 2013 Robert Hülle
// No warranty. WTFPL v2

/*

Hashmv appends hash to filename.

example:
	touch file foo.ext
	hashmv file
output is files named
	file[00000000]
	foo[00000000].ext

Files renamed by this program should be accepted by hashcheck.

*/
package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

var (
	pretend bool
)

func init() {
	flag.BoolVar(&pretend, "p", false, "do not rename file")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	for _, file := range flag.Args() {
		fileHash, err := hashCrc32(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(10)
		}
		name := newName(file, fileHash)
		if !pretend {
			err = os.Rename(file, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
				os.Exit(11)
			}
		}
		fmt.Printf("`%s` renamed as: `%s`\n", file, name)
	}
}

func hashCrc32(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buff := make([]byte, 8192)
	hash := crc32.NewIEEE()
	var length int
	for length, err = file.Read(buff); length > 0; length, err = file.Read(buff) {
		hash.Write(buff[:length])
	}
	if err != nil && err != io.EOF {
		return fmt.Sprintf("%08X", hash.Sum32()), err
	}
	return fmt.Sprintf("%08X", hash.Sum32()), nil
}

func newName(file, fileHash string) string {
	dot := strings.LastIndex(file, ".")
	if dot == -1 {
		dot = len(file)
	}
	name := file[0:dot] + "[" + fileHash + "]" + file[dot:len(file)]
	return name
}

func usage() {
	fmt.Fprintln(os.Stderr, "hashmv - append hash to filename\n")
	fmt.Fprintf(os.Stderr, "usage: %s [options] <file> ...\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "options:\n")
	flag.PrintDefaults()
}
