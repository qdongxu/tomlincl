package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/qdongxu/tomlincl/incl"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s <toml file>", os.Args[0])
	}
	if len(os.Args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	buf, err := incl.ParseIncludeRecursively(os.Args[1], &bytes.Buffer{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err())
		os.Exit(3)
	}
}
