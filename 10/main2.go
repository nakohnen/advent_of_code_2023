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

// Contains checks if a slice contains a particular element
func contains[T comparable](s []T, elem T) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
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
	filename = ""
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
	lines = nil
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
	new_connection = nil
	// Now we wander the map
	to_go := [][2]int{start}
	visited := make(map[[2]int]bool)
	step := -1
	for len(to_go) > 0 {
		//fmt.Printf("%v step=%v\n", to_go, step)
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
	loop := visited
	fmt.Println("Loop found.")

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pos := [2]int{y, x}
			if !loop[pos] {
				tile_map[y][x] = '.'
			}
		}
	}
	fmt.Println("Tilemap cleaned.")

	down_pointing := map[rune]bool{'|': true, 'F': true, '7': true, 'S': true}
	up_pointing := map[rune]bool{'|': true, 'J': true, 'L': true, 'S': true}
	right_pointing := map[rune]bool{'-': true, 'F': true, 'L': true, 'S': true}
	left_pointing := map[rune]bool{'-': true, 'J': true, '7': true, 'S': true}

	// Now we double the size and connect the pipes to make "room" to breathe.
	tile_map2 := [][]rune{}
	loop2 := make(map[[2]int]bool)
	for y := 0; y < height*2; y++ {
		new_row := []rune{}
		for x := 0; x < width*2; x++ {
			new_tile := '.'
			if y == height*2-1 || x == width*2-1 {
				new_tile = '.'
			} else if y%2 == 0 && x%2 == 0 {
				new_tile = tile_map[y/2][x/2]
			} else if y%2 == 1 && x%2 == 1 {
				new_tile = '.'
			} else if y%2 == 1 {
				upper_tile := tile_map[(y-1)/2][x/2]
				lower_tile := tile_map[(y+1)/2][x/2]
				if down_pointing[upper_tile] && up_pointing[lower_tile] {
					new_tile = '|'
				}
			} else if x%2 == 1 {
				left_tile := tile_map[y/2][(x-1)/2]
				right_tile := tile_map[y/2][(x+1)/2]
				if right_pointing[left_tile] && left_pointing[right_tile] {
					new_tile = '-'
				}
			}
			new_row = append(new_row, new_tile)
		}
		tile_map2 = append(tile_map2, new_row)
	}

	for y := 0; y < height*2; y++ {
		for x := 0; x < width*2; x++ {
			pos := [2]int{y, x}
			if tile_map2[y][x] != '.' {
				loop2[pos] = true
			}
			//fmt.Printf("%v", string(tile_map2[y][x]))
		}
		//fmt.Println("")
	}


	start = [2]int{0, 0}
	to_go = [][2]int{start}
	directions := [2]int{-1, 1}
    visited = make(map[[2]int]bool)

    //break_point := 1000
	for len(to_go) > 0 {
		pos := to_go[0]
		to_go = to_go[1:]
		visited[pos] = true

		for _, d := range directions {
			if 0 <= pos[0]+d && pos[0]+d < height*2 {
				new_pos := [2]int{pos[0] + d, pos[1]}
				if !visited[new_pos] && !loop2[new_pos] {
                    visited[new_pos] = true
					to_go = append(to_go, new_pos)
				}
			}
			if 0 <= pos[1]+d && pos[1]+d < width*2 {
				new_pos := [2]int{pos[0], pos[1] + d}
				if !visited[new_pos] && !loop2[new_pos] {
                    visited[new_pos] = true
					to_go = append(to_go, new_pos)
				}
			}
		}
        //if len(to_go) > break_point {
        //    break
        //}
        //fmt.Printf("%v len(to_go) %v\n", len(to_go), to_go)
	}
	outside2 := visited
	visited = make(map[[2]int]bool)
	fmt.Println("First positions of outside marked.")

    sum := 0
    for y:=0;y<height;y++ {
        for x:=0;x<width;x++{
            pos := [2]int{y, x}
            pos2 := [2]int{y*2, x*2}
            if !loop[pos] && !outside2[pos2] {
                sum++
            }
        }
    }

	fmt.Printf("sum=%v\n", sum)
	fmt.Println("")
}
