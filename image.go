package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type Image struct {
	width  int
	height int
	pixels []uint8
}

func ReadImgFromFile(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	return string(data)
}

func LoadDirectory(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.EqualFold(filepath.Ext(entry.Name()), ".txt") {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}

	return files, nil
}

func OldParseImage(data string, width, height int) Image {
	fields := strings.Fields(data)

	img := Image{
		width:  width,
		height: height,
		pixels: make([]uint8, len(fields)),
	}

	for i, field := range fields {
		n, err := strconv.Atoi(field)
		if err != nil {
			panic(err)
		}

		img.pixels[i] = uint8(n)
	}

	return img
}

func ParseImage(data string) Image {
	// Find the end of the first line
	lineEnd := strings.Index(data, "\n")

	// Extract the first line
	firstLine := data[:lineEnd]

	// Split into two values
	dimensions := strings.Fields(firstLine)

	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	// Remove the first line from the original string
	data = data[lineEnd+1:]

	data = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, data)

	if len(data) != width*height {
		fmt.Printf("Length: %d (expected %d)\n", len(data), width*height)
		panic("incorrect image data length")
	}

	img := Image{
		width:  width,
		height: height,
		pixels: make([]uint8, len(data)),
	}

	for i := 0; i < len(data); i++ {
		if data[i] != '0' && data[i] != '1' {
			panic(fmt.Sprintf("invalid pixel at position %d: %q", i, data[i]))
		}

		img.pixels[i] = data[i] - '0'
	}

	return img
}
