package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const debug bool = false

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

type Vector struct {
	pos    Position
	norm   Position
	length int
}

func getNorm(in string) Position {
	switch in {
	case "R":
		return Position{x: 1, y: 0}
	case "L":
		return Position{x: -1, y: 0}
	case "U":
		return Position{x: 0, y: -1}
	case "D":
		return Position{x: 0, y: 1}
	}
	return Position{x: 0, y: 0}
}

func multiplyPos(pos Position, scalar int) Position {
	return Position{pos.x * scalar, pos.y * scalar}
}

func addPos(a, b Position) Position {
	return Position{a.x + b.x, a.y + b.y}
}

func createPartitions(figure []Vector) [][2]int {
	result := [][2]int{}
	break_points := []int{}
	for _, v := range figure {
		break_points = append(break_points, v.pos.y)
		break_points = append(break_points, v.pos.y+v.norm.y*v.length)
	}
	break_points = removeDuplicates(break_points)
	sort.Ints(break_points)
	for i := 1; i < len(break_points); i++ {
		result = append(result, [2]int{break_points[i-1], break_points[i-1]})
		result = append(result, [2]int{break_points[i-1] + 1, break_points[i] - 1})
	}
	last := break_points[len(break_points)-1]
	result = append(result, [2]int{last, last})
	new_result := [][2]int{}
	for _, t := range result {
		if t[0] <= t[1] {
			new_result = append(new_result, t)
		}
	}
	return new_result
}

func getEndpoints(vec Vector) (Position, Position) {
    v1 := vec.pos
    v2 := addPos(vec.pos, multiplyPos(vec.norm, vec.length))
    if v1.x < v2.x || (v1.x == v2.x && v1.y < v2.y ){
        return v1, v2
    }
    return v2, v1
}

func calculateArea(partition [2]int, figure []Vector) int {
	y := partition[0]
	result := 0

	v_vectors := []Vector{}
	h_vectors := []Vector{}
	for _, v := range figure {
		start := v.pos
		end := addPos(v.pos, multiplyPos(v.norm, v.length))
		if min(start.y, end.y) <= y && y < max(start.y, end.y) && abs(v.norm.y) == 1 {
				if debug {
					fmt.Printf("y=%v -> v=%v (%v to %v) v", y, v, start, end)
					fmt.Print(" appended.")
					fmt.Println("")
				}
				v_vectors = append(v_vectors, v)
		} else if abs(v.norm.x) == 1 && v.pos.y == y{
				if debug {
					fmt.Printf("y=%v -> v=%v (%v to %v) h", y, v, start, end)
					fmt.Print(" appended.")
					fmt.Println("")
				}
				h_vectors = append(h_vectors, v)
		} /*else {
            if debug {
                fmt.Printf("Dropped %v\n", v)
            }
        } */
	}

	sort.Slice(v_vectors, func(i, j int) bool { return v_vectors[i].pos.x < v_vectors[j].pos.x })
	sort.Slice(h_vectors, func(i, j int) bool { return h_vectors[i].pos.x < h_vectors[j].pos.x })
    covered := [][2]int{}
	for i := 0; i < len(v_vectors)/2; i++ {
		v1 := v_vectors[i*2]
		v2 := v_vectors[i*2+1]
        covered = append(covered, [2]int{v1.pos.x, v2.pos.x} )
        r := abs(v1.pos.x - v2.pos.x) + 1
        if debug {
            fmt.Printf("For partition %v we have the pair %v, %v with width %v\n",partition, v1, v2, r)
        }
        result += r
	}
    for _, h_v := range h_vectors {
        if debug {
            fmt.Printf("Lets verify %v for extra.\n", h_v)
        }
        found := false
        start, end := getEndpoints(h_v)
        for _, interval := range covered {
           if interval[0] <= start.x && end.x <= interval[1]  {
                found = true
                break
           }
        }
        if !found{
            r := h_v.length
            result += r
            if debug {
                fmt.Printf("We got the additional %v with width %v\n", h_v, r)
            }
        }
    }

	return result * (partition[1] - partition[0] + 1)
}

func processWorker(id int, work <-chan [2]int, resultChan chan<- int, wg *sync.WaitGroup,
	figure []Vector, width int) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		//result := processMirrorImage(w, false)
		result := calculateArea(w, figure)
		resultChan <- result
		fmt.Printf("Worker=%v: %v => %v\n", id, w, result)
	}
}

func removeDuplicates(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, i := range slice {
		if !seen[i] {
			result = append(result, i)
			seen[i] = true
		}
	}
	return result
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
	vectors := []Vector{}
	x_0 := 0
	y_0 := 0
	height := 0
	width := 0
	position := Position{x: x_0, y: y_0}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			fmt.Printf("%v ", cleaned_line)

			line_split := strings.Split(cleaned_line, " ")
			hexa := line_split[2]
			dir := ""
			switch dir_h := string(hexa[7]); dir_h {
			case "0":
				dir = "R"
			case "1":
				dir = "D"
			case "2":
				dir = "L"
			case "3":
				dir = "U"
			}
			if dir != "R" && dir != "L" && dir != "U" && dir != "D" {
				panic("Error parsing")
			}
			length, err := strconv.ParseInt(string(hexa[2:7]), 16, 64)
			check(err)

			norm := getNorm(dir)
			vectors = append(vectors, Vector{position, norm, int(length)})

			position = addPos(position, multiplyPos(norm, int(length)))
			fmt.Printf("=> %v %v  %v (%v)\n", dir, length, hexa, position)
			x_0 = min(position.x, x_0)
			y_0 = min(position.y, y_0)
			width = max(position.x, width)
			height = max(position.y, height)
		}
	}
	fmt.Printf("x_0=%v y_0=%v, w=%v h=%v\n", x_0, y_0, width, height)
	width = width - x_0 + 2
	height = height - y_0 + 2

	for i := 0; i < len(vectors); i++ {
		v := vectors[i]
		v.pos.x -= x_0 - 1
		v.pos.y -= y_0 - 1
		vectors[i] = v
	}

	if debug {
		for _, v := range vectors {
			if v.pos.x < 0 || v.pos.x >= width {
				panic("Error in the position translation.")
			}
			if v.pos.y < 0 || v.pos.y >= height {
				panic("Error in the position translation.")
			}
			fmt.Printf("%v\n", v)
		}
		fmt.Printf("x_0=%v y_0=%v, w=%v h=%v\n", x_0, y_0, width, height)
	}

	//**************************+
	partitions := createPartitions(vectors)
	channel_length := len(partitions)
	resultChan := make(chan int, channel_length)
	workChan := make(chan [2]int, channel_length)
	var wg sync.WaitGroup

	var numWorkers int = cores
	if debug {
		numWorkers = 1
	} else {
		if numWorkers < 1 {
			numWorkers = 1
		} else if numWorkers > 24 {
			numWorkers = 24
		}
		numWorkers = min(numWorkers, channel_length)
	}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processWorker(i, workChan, resultChan, &wg, vectors, width)
	}
	fmt.Printf("To work on: %v starts with %v threads\n", channel_length, numWorkers)

	for i := 0; i < channel_length; i++ {
		workChan <- partitions[i]
	}
	close(workChan)

	// Close the channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := 1
	for result := range resultChan {
		results += result
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
