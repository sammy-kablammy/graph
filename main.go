package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x float64
	y float64
}

func main() {
	// parse command-line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Draw and ASCII graph based on the coordinates from stdin.\n")
		fmt.Fprintf(os.Stderr, "Coordinates should be one pair per line, separated by white space.\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "\t1 10\n")
		fmt.Fprintf(os.Stderr, "\t-2 -2\n")
		fmt.Fprintf(os.Stderr, "\t5.1 -3\n")
		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	var termColumns int
	var termLines int
	flag.IntVar(&termColumns, "w", 80, "width of graph in characters")
	flag.IntVar(&termLines, "l", 24, "height of graph in lines")
	flag.Parse()

	// read input
	scan := bufio.NewScanner(os.Stdin)
	points := []point{}
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	minY := math.Inf(1)
	maxY := math.Inf(-1)
	for scan.Scan() {
		line := scan.Text()
		split := strings.Fields(line)
		if len(split) != 2 {
			fmt.Fprintf(os.Stderr, "Expected exactly two values per row\n")
			os.Exit(1)
		}

		x, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid x-coordinate '%s'\n", split[0])
			os.Exit(1)
		}
		y, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid y-coordinate '%s'\n", split[1])
			os.Exit(1)
		}

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}

		points = append(points, point{x, y})
	}

	// map the "real" coordinates onto screen coordinates
	newPoints := []point{}
	for _, p := range points {
		newX := mapRange(minX, maxX, 0, float64(termColumns), p.x)
		newY := mapRange(minY, maxY, 0, float64(termLines), p.y)
		newPoints = append(newPoints, point{newX, newY})
	}

	// initialize output
	output := make([][]byte, termLines+1)
	for i := range output {
		output[i] = make([]byte, termColumns+1)
		for charIdx := range output[i] {
			output[i][charIdx] = '.'
		}
	}


	// plot some points (the thing you wanted to do the whole time)
	const plottedPoint = '#'
	for _, p := range newPoints {
		output[int(p.y)][int(p.x)] = plottedPoint
	}

	printOutputWithBorder(output)
}

// Map the value 't' from the range [inStart, inEnd] to [outStart, outEnd].
func mapRange(inStart, inEnd, outStart, outEnd, t float64) float64 {
	return outStart + ((outEnd-outStart)/(inEnd-inStart))*(t-inStart)
}

// https://en.wikipedia.org/wiki/Box-drawing_characters
const borderTop = "━"
const borderSide = "┃"
const borderCornerNE = "┓"
const borderCornerNW = "┏"
const borderCornerSE = "┛"
const borderCornerSW = "┗"

func printOutputWithBorder(output [][]byte) {
	// top border
	fmt.Print(borderCornerNW)
	for range len(output[0]) {
		fmt.Print(borderTop)
	}
	fmt.Print(borderCornerNE)
	fmt.Println()
	// left border, body, right border
	for i := len(output) - 1; i >= 0; i-- {
		fmt.Print(borderSide)
		fmt.Print(string(output[i]))
		fmt.Print(borderSide)
		fmt.Println()
	}
	// bottom border
	fmt.Print(borderCornerSW)
	for range len(output[0]) {
		fmt.Print(borderTop)
	}
	fmt.Print(borderCornerSE)
	fmt.Println()
}
