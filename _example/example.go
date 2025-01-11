package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/qdongxu/tomlincl/incl"

	"github.com/BurntSushi/toml"
)

type (
	example struct {
		Title      string
		Desc       string
		Integers   []int
		Floats     []float64
		Times      []fmtTime
		Duration   []time.Duration
		Distros    []distro
		Servers    map[string]server
		Characters map[string][]struct {
			Name string
			Rank string
		}
		Services []service
	}

	server struct {
		IP       string
		Hostname string
		Enabled  bool
	}

	distro struct {
		Name     string
		Packages string
	}

	service struct {
		Name  string
		Nodes []node
	}

	node struct {
		IP string
	}
	fmtTime struct{ time.Time }
)

func (t fmtTime) String() string {
	f := "2006-01-02 15:04:05.999999999"
	if t.Time.Hour() == 0 {
		f = "2006-01-02"
	}
	if t.Time.Year() == 0 {
		f = "15:04:05.999999999"
	}
	if t.Time.Location() == time.UTC {
		f += " UTC"
	} else {
		f += " -0700"
	}
	return t.Time.Format(`"` + f + `"`)
}

func main() {
	f := "example.toml"
	if _, err := os.Stat(f); err != nil {
		f = "_example/example.toml"
	}

	buf, err := incl.ParseIncludeRecursively(f, &bytes.Buffer{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data, err := io.ReadAll(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var config example
	meta, err := toml.Decode(string(data), &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(config.Services) != 2 {
		fmt.Fprintln(os.Stderr, "failed to include the first level")
		os.Exit(1)
	}

	if len(config.Services[0].Nodes) != 2 {
		fmt.Fprintln(os.Stderr, "failed to include the second level")
		os.Exit(1)
	}

	indent := strings.Repeat(" ", 14)

	fmt.Print("Decoded")
	typ, val := reflect.TypeOf(config), reflect.ValueOf(config)
	for i := 0; i < typ.NumField(); i++ {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 7)
		}
		fmt.Printf("%s%-11s → %v\n", indent, typ.Field(i).Name, val.Field(i).Interface())
	}

	fmt.Print("\nKeys")
	keys := meta.Keys()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	for i, k := range keys {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 10)
		}
		fmt.Printf("%s%-10s %s\n", indent, meta.Type(k...), k)
	}

	fmt.Print("\nUndecoded")
	keys = meta.Undecoded()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	for i, k := range keys {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 5)
		}
		fmt.Printf("%s%-10s %s\n", indent, meta.Type(k...), k)
	}
}
