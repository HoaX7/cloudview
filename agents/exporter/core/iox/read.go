package iox

import (
	"bufio"
	"os"
	"strings"
)

type (
	textReadOptions struct {
		keepSpace     bool
		withoutBlanks bool
		omitPrefix    string
	}

	// TextReadOption defines the method to customize the text reading functions.
	TextReadOption func(*textReadOptions)
)

// KeepSpace customizes the reading functions to keep leading and tailing spaces.
func KeepSpace() TextReadOption {
	return func(o *textReadOptions) {
		o.keepSpace = true
	}
}

// ReadText reads content from the given file with leading and tailing spaces trimmed.
func ReadText(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// ReadTextLines reads the text lines from given file.
func ReadTextLines(filename string, opts ...TextReadOption) ([]string, error) {
	var readOpts textReadOptions
	for _, opt := range opts {
		opt(&readOpts)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !readOpts.keepSpace {
			line = strings.TrimSpace(line)
		}
		if readOpts.withoutBlanks && len(line) == 0 {
			continue
		}
		if len(readOpts.omitPrefix) > 0 && strings.HasPrefix(line, readOpts.omitPrefix) {
			continue
		}

		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

// WithoutBlank customizes the reading functions to ignore blank lines.
func WithoutBlank() TextReadOption {
	return func(o *textReadOptions) {
		o.withoutBlanks = true
	}
}

// OmitWithPrefix customizes the reading functions to ignore the lines with given leading prefix.
func OmitWithPrefix(prefix string) TextReadOption {
	return func(o *textReadOptions) {
		o.omitPrefix = prefix
	}
}
