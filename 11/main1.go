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

func manhattanMetric(pos1 [2]int, pos2 [2]int) int {
	return abs(pos1[0]-pos2[0]) + abs(pos1[1]-pos2[1])
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

	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	height := len(lines)
	width := len(lines[0])
	systems1 := [][2]int{}
	expansions_v := make(map[int]bool)
	expansions_h := make(map[int]bool)
	for y, line := range lines {
		for x, r := range line {
			pos := [2]int{y, x}

			if r == '#' {
				systems1 = append(systems1, pos)
				expansions_v[x] = true
				expansions_h[y] = true
			}
		}
	}
	y_expansion := 0
	systems2 := [][2]int{}
	for y := 0; y < height; y++ {
		if !expansions_h[y] {
			y_expansion++
		}
        x_expansion := 0
		for x := 0; x < width; x++ {
			if !expansions_v[x] {
				x_expansion++
			}

			for _, s := range systems1 {
				if y == s[0] && x == s[1] {
					new_pos := [2]int{y_expansion + s[0], x_expansion + s[1]}
					systems2 = append(systems2, new_pos)
				}
			}

		}
	}
	done := make(map[[2]int]bool)
	sum := 0
	for id1, s1 := range systems2 {
		for id2, s2 := range systems2 {
			ma := max(id1, id2)
			mi := min(id1, id2)
			if !done[[2]int{ma, mi}] {
				d := manhattanMetric(s1, s2)
				sum += d
				done[[2]int{ma, mi}] = true
				fmt.Printf("%v %v d=%v\n", s1, s2, d)
			}
		}
	}
	fmt.Printf("sum=%v\n", sum)
	fmt.Println("")
}
