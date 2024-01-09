package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PBM est une structure pour représenter des images PBM.
type PBM struct {
	Data          [][]bool
	Width, Height int
	MagicNumber   string
}

// ReadPBM lit une image PBM à partir d'un fichier et renvoie une structure représentant l'image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read magic number")
	}
	magicNumber := scanner.Text()

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("unsupported PBM format: %s", magicNumber)
	}

	// Skip comments and empty lines
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' {
			break
		}
	}

	// Read width and height
	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading dimensions line: %v", scanner.Err())
	}
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("invalid dimensions line")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse width: %v", err)
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse height: %v", err)
	}

	// Read data
	var data [][]bool
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		row := make([]bool, width)
		for i, token := range tokens {
			if i >= width {
				break
			}
			if token == "1" {
				row[i] = true
			} else if token == "0" {
				row[i] = false
			} else {
				return nil, fmt.Errorf("invalid character in data: %s", token)
			}
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return &PBM{
		Data:        data,
		Width:       width,
		Height:      height,
		MagicNumber: magicNumber,
	}, nil
}