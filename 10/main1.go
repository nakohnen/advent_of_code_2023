package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
	//	"math"
	//"sort"
	"strconv"
	//"sync"
)

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

// allElementsSame checks if all elements of the slice are the same
func allElementsSame(slice []int) bool {
	if len(slice) == 0 {
		return true // Optionally, define behavior for empty slices
	}

	firstElement := slice[0]
	for _, element := range slice {
		if element != firstElement {
			return false
		}
	}
	return true
}

func main() {
	// Check if an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

	// The second element in os.Args is the first argument
	filename := os.Args[1]

	dat, err := os.ReadFile(filename)
	check(err)

	tile_map := [][]rune{}
	start := [2]int{-1, -1}
	connection := make(map[[2]int][][2]int)
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	height := len(lines)
	width := len(lines[0])
	for y, line := range lines {
		new_row := []rune{}
		for x, r := range line {
			new_row = append(new_row, r)
			pos := [2]int{y, x}
			connection[pos] = [][2]int{}

			left := false
			right := false
			up := false
			down := false

			switch r {
			case 'S':
				start = pos
				left = true
				right = true
				up = true
				down = true
			case '|':
				up = true
				down = true
			case '-':
				left = true
				right = true
			case '7':
				left = true
				down = true
			case 'J':
				up = true
				left = true
			case 'F':
				right = true
				down = true
			case 'L':
				up = true
				right = true
			}

			if ty := y - 1; up && ty >= 0 {
				connection[pos] = append(connection[pos], [2]int{ty, x})
			}
			if ty := y + 1; down && ty < height {
				connection[pos] = append(connection[pos], [2]int{ty, x})
			}
			if tx := x - 1; left && tx >= 0 {
				connection[pos] = append(connection[pos], [2]int{y, tx})
			}
			if tx := x + 1; right && tx < width {
				connection[pos] = append(connection[pos], [2]int{y, tx})
			}
		}
		tile_map = append(tile_map, new_row)
	}
    //fmt.Printf("connection_map: \n%v\n", connection)
	new_connection := make(map[[2]int][][2]int)
	// We scanned all the lines and created a rought connection map. Now we need to trim the map.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pos := [2]int{y, x}
			new_connection[pos] = [][2]int{}
			for _, t_pos := range connection[pos] {
				found := false
				for _, b_pos := range connection[t_pos] {
      //              fmt.Printf("%v==%v?\n", b_pos, pos)
					if b_pos == pos {
						found = true
						break
					}
				}

				if found {
    //                fmt.Printf("%v and %v are connected.\n", pos, t_pos )
					new_connection[pos] = append(new_connection[pos], t_pos)
				}
			}
		}
	}
	connection = new_connection
  //  fmt.Printf("connection_map: \n%v\n", connection)
	new_connection = make(map[[2]int][][2]int)
	// Now we wander the map
	to_go := [][2]int{start}
	visited := make(map[[2]int]bool)
	step := -1
	for len(to_go) > 0 {
		fmt.Printf("%v step=%v\n", to_go, step)
		step++
		new_to_go := [][2]int{}
		for _, pos := range to_go {
			visited[pos] = true
			for _, t_pos := range connection[pos] {
//                fmt.Printf("Visiting %v\n", t_pos)
				if !visited[t_pos] {
					new_to_go = append(new_to_go, t_pos)
//					fmt.Printf("%v -> %v\n", pos, t_pos)
				}
			}
		}
		to_go = new_to_go
	}

	fmt.Printf("step=%v\n", step)
	fmt.Println("")
}
