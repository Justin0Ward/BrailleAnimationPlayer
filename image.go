package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func ParseImage(data string, width, height int) Image {
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
