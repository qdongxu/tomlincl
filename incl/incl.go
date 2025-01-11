package incl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var includeDirective = regexp.MustCompile(`^\s*#!\s*include\s*(.*)\s*$`)

func ParseIncludeRecursively(path string, buf *bytes.Buffer) (io.Reader, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		match := includeDirective.FindStringSubmatch(line)
		if len(match) == 0 {
			buf.WriteString(line)
			buf.WriteByte('\n')
			continue
		} else {
			fmt.Println("abc   " + match[1])
		}

		cwd := filepath.Dir(path)
		includes, err := filepath.Glob(filepath.Join(cwd, match[1]))
		if err != nil {
			return nil, err
		}

		for _, nextFile := range includes {
			_, err = ParseIncludeRecursively(nextFile, buf)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s, error: %w", nextFile, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan %s, error: %w", path, err)
	}

	return buf, nil
}
