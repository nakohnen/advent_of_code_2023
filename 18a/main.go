package main

import (
	"flag"
	"fmt"
	//"math"
	"os"
	//"sort"
	"strconv"
	"strings"
	"sync"
	//"unicode"
)

const debug bool = true

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toInt(val string) int {
	v, err := strconv.Atoi(val)
	check(err)
	return v
}

func toRune(val int) rune {
	return rune('0' + val)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func hasPrefix(line string, sub string) bool {
	if len(line) < len(sub) {
		return false
	}
	return line[0:len(sub)] == sub
}

func replaceCharacters(line string, rem string, rep string) string {
	result := line
	for _, c := range rem {
		result = strings.ReplaceAll(result, string(c), rep)
	}
	return result
}

func sumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
}

func IndexOf[T comparable](s []T, element T) int {
	for i, e := range s {
		if e == element {
			return i
		}
	}
	return -1
}

func IndexOfString(s string, r rune) int {
	for i, r2 := range s {
		if r2 == r {
			return i
		}
	}
	return -1
}

type Position struct {
	x, y int
}

func processWorker(id int, work <-chan Position, resultChan chan<- int, wg *sync.WaitGroup,
	grid map[Position]int, width, height int, target Position) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		//result := processMirrorImage(w, false)
		result := 0
		resultChan <- 0
		fmt.Printf("Worker=%v: %v => %v\n", id, w, result)
	}
}

func getSteps(x, y int, direction string, steps int) []Position {
	result := []Position{}
	switch direction {
	case "R":
		for sx := x; sx <= x+steps; sx++ {
			result = append(result, Position{x: sx, y: y})
		}
	case "L":
		for sx := x; sx >= x-steps; sx-- {
			result = append(result, Position{x: sx, y: y})
		}
	case "D":
		for sy := y; sy <= y+steps; sy++ {
			result = append(result, Position{x: x, y: sy})
		}
	case "U":
		for sy := y; sy >= y-steps; sy-- {
			result = append(result, Position{x: x, y: sy})
		}
	}
	if debug {
		fmt.Printf("%v\n", result)
	}

	return result
}

func drawGrid() {

}

func main() {
	if debug {
		fmt.Println("Debug enabled")
	}
	cores_f := flag.Int("t", -1, "How many cores (threads) should we run?")
	filename_f := flag.String("f", "", "On which file should we run this?")
	flag.Parse()

	// The second element in os.Args is the first argument
	filename := *filename_f
	if len(filename) == 0 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}
	cores := *cores_f
	if cores < 0 {
		cores = 1
	}

	dat, err := os.ReadFile(filename)
	check(err)

	// Read file
	grid_color := make(map[Position]string)
	positions := []Position{}
	x_0 := 0
	y_0 := 0
	height := 0
	width := 0
	position := Position{x: x_0, y: y_0}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			fmt.Printf("%v\n", cleaned_line)

			line_split := strings.Split(cleaned_line, " ")
			dir := line_split[0]
			length := toInt(line_split[1])
			color := line_split[2]
			for _, new_pos := range getSteps(position.x, position.y, dir, length) {
				grid_color[new_pos] = color
				x_0 = min(x_0, new_pos.x)
				y_0 = min(y_0, new_pos.y)
				height = max(height, new_pos.y)
				width = max(width, new_pos.x)
				positions = append(positions, new_pos)
				position = new_pos
			}
		}
	}
	fmt.Printf("x_0=%v y_0=%v, w=%v h=%v\n", x_0, y_0, width, height)
	width = width - x_0 + 2
	height = height - y_0 + 2
	grid_color_new := make(map[Position]string)
	border := make(map[Position]bool)
	border_count := make(map[Position]int)
	for x := 0; x <= width; x++ {
		for y := 0; y <= height; y++ {
			border_count[Position{x, y}] = 0
		}
	}
	for i := 0; i < len(positions); i++ {
		color := grid_color[positions[i]]
		positions[i].x -= x_0 - 1
		positions[i].y -= y_0 - 1
		grid_color_new[positions[i]] = color
		border[positions[i]] = true
		border_count[positions[i]] += 1
		if debug {
			if border_count[positions[i]] > 1 {
				fmt.Printf("Corner at x=%v y=%v\n", positions[i].x, positions[i].y)
			}
		}
	}
    fmt.Println("")
    work := []Position{{0, 0}}
    
    outside := make(map[Position]bool)
    outside[Position{0,0}] = true
    for len(work) > 0 {
        pos := work[0]
        work = work[1:]
        

        new_candidates := []Position{}
        new_candidates = append(new_candidates, Position{x: pos.x-1,y:pos.y})
        new_candidates = append(new_candidates, Position{x: pos.x+1,y:pos.y})
        new_candidates = append(new_candidates, Position{x: pos.x,y:pos.y-1})
        new_candidates = append(new_candidates, Position{x: pos.x,y:pos.y+1})
    
        for _, cand := range new_candidates {
            if cand.x >= 0 && cand.y >=0 && cand.x <= width && cand.y <= height {
                if !border[cand] && !outside[cand] {
                    work = append(work, cand)
                    outside[cand] = true
                }
            }
        }
    }
	fmt.Println("")
	grid_color = grid_color_new
	grid_color_new = nil
	grid := [][]int{}
	var sb strings.Builder
    results := 0
	for y := 0; y <= height; y++ {
		sb.Reset()
		line := []int{}
		for x := 0; x <= width; x++ {
			pos := Position{x: x, y: y}
			val := 0
			if border_count[pos] > 0 {
				val = 1
			} else if !outside[pos] {
                val = 2
            }
			line = append(line, val)
			sb.WriteRune(toRune(val))
            if val > 0 {
                results++
            }
		}
		grid = append(grid, line)
        if debug {
		    fmt.Printf("%v\n", sb.String())
        }
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
