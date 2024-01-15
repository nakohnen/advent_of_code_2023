package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	//"unicode"
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

func trackBack(vec Vector, came_from map[Vector]Position) map[Position]bool {
	result := make(map[Position]bool)
	result[vec.pos] = true
	next, ok := came_from[vec]
	new_direction := HORIZONTAL
	if vec.dir == HORIZONTAL {
		new_direction = VERTICAL
	}
	for ok {
		result[next] = true
		if new_direction == HORIZONTAL {
			new_direction = VERTICAL
		} else {
			new_direction = HORIZONTAL
		}
		next, ok = came_from[Vector{next, new_direction}]
	}

	return result
}

func paintGrid(grid map[Position]int, height, width int, came_from map[Vector]Position, current WalkCost) {
	var sb strings.Builder
	var sb_right strings.Builder
	current_vec := Vector{current.pos, current.direction}
	track := trackBack(current_vec, came_from)
	for y := 0; y < height; y++ {
		sb.Reset()
		sb_right.Reset()
		sb_right.WriteString(" ")
		for x := 0; x < width; x++ {
			pos := Position{x, y}
			var tile string = fmt.Sprintf("%v", grid[pos])
			var right_tile string = "."
			if track[pos] {
				tile = "#"
				right_tile = fmt.Sprintf("%v", grid[pos])
			}

			sb.WriteString(tile)
			sb_right.WriteString(right_tile)
		}
		fmt.Println(sb.String() + sb_right.String())
	}
}

type Direction int

const (
	HORIZONTAL Direction = iota
	VERTICAL
)

func orthogonal(a, b Direction) bool {
	if a != b {
		return true
	}
	return false
}

type WalkCost struct {
	pos       Position
	cost      int
	direction Direction
    f_value   float64
}

// copyMap creates a new map and copies the key-value pairs from the original map
func copyMap(originalMap map[Position]bool) map[Position]bool {
	newMap := make(map[Position]bool)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

type Vector struct {
	pos Position
	dir Direction
}

func cleanOpenSet(set []WalkCost, target Position) []WalkCost {
    result := []WalkCost{}
    is_present := make(map[WalkCost]bool)
    for _, wc := range set {
        found := false
        min_cost := math.MaxInt
        for _, r_wc := range result {
            if wc.pos == r_wc.pos && wc.direction == r_wc.direction {
                found = true
                if wc.cost < min_cost {
                    min_cost = wc.cost
                }
            }
        }
        if !found {
            result = append(result, wc)
        } else if !is_present[wc] {
            result = append(result, WalkCost{wc.pos, min_cost, wc.direction, calculateFValue(wc.pos, target, min_cost)})
            is_present[wc] = true
        }
    }
    return result
}

func calculateFValue(source, target Position, cost int) float64 {
            d_b2 := math.Pow(float64(source.x - target.x), 2) + math.Pow(float64(source.y - target.y),2)
            d_a := math.Sqrt(d_b2)
            return d_a +  float64(cost)
}

func walkGrid(grid map[Position]int, width, height int, start Position, target Position) int {
	open_set := []WalkCost{}
	open_set = append(open_set, WalkCost{start, 0, HORIZONTAL, calculateFValue(start, target, 0)})
	open_set = append(open_set, WalkCost{start, 0, VERTICAL, calculateFValue(start, target, 0)})
	cost_map := make(map[Vector]int)
	came_from := make(map[Vector]Position)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for _, dir := range []Direction{HORIZONTAL, VERTICAL} {
				cost_map[Vector{Position{x, y}, dir}] = math.MaxInt
			}
		}
	}
	for _, wc := range open_set {
		cost_map[Vector{wc.pos, wc.direction}] = 0
	}
	current := open_set[0]

	for len(open_set) > 0 {
        if len(open_set) > 5000 {
            open_set = cleanOpenSet(open_set, target)
        }
		sort.Slice(open_set, func(i, j int) bool {
            return open_set[i].f_value < open_set[j].f_value
		})
		if debug {
            if len(open_set) > 6 {
                fmt.Printf("Open Set: %v\n", len(open_set))
            } else {
                fmt.Printf("Open Set: %v\n", open_set)
            }
		}
		current = open_set[0]

		open_set = open_set[1:]
		if debug {

			fmt.Printf("Current: %v ", current)
			val, ok := came_from[Vector{current.pos, current.direction}]
			if ok {
				fmt.Printf("came from %v", val)
			}
			fmt.Println("")
			//paintGrid(grid, height, width, came_from, current)
		}

		if current.pos == target {
			return current.cost
		}

		candidates := []Vector{}
            for i := 4 ; i <= 10 ; i ++ {
                if current.direction == VERTICAL {
                candidates = append(candidates, Vector{Position{current.pos.x - i, current.pos.y}, HORIZONTAL})
                candidates = append(candidates, Vector{Position{current.pos.x + i, current.pos.y}, HORIZONTAL})
            } else if current.direction == HORIZONTAL {
                candidates = append(candidates, Vector{Position{current.pos.x, current.pos.y - i}, VERTICAL})
                candidates = append(candidates, Vector{Position{current.pos.x, current.pos.y + i}, VERTICAL})
            }
		}
		for _, c := range candidates {
			pos := c.pos
			if debug {
				fmt.Printf("%v candidate.\n", c)
			}
			if pos.x >= 0 && pos.y >= 0 && pos.x < width && pos.y < height {
				grid_cost := current.cost
				switch dir := c.dir; dir {
				case HORIZONTAL:
					for x := min(pos.x, current.pos.x); x <= max(pos.x, current.pos.x); x++ {
                        if x != current.pos.x {
						    grid_cost += grid[Position{x, current.pos.y}]
                        }
					}

				case VERTICAL:
					for y := min(pos.y, current.pos.y); y <= max(pos.y, current.pos.y); y++ {
                        if y != current.pos.y {
    						grid_cost += grid[Position{current.pos.x, y}]
				    	}
                    }
				}
				new_cost := grid_cost
				if new_cost <= cost_map[c] {
					wc := WalkCost{c.pos, new_cost, c.dir, calculateFValue(c.pos, target, new_cost)}
					open_set = append(open_set, wc)
					cost_map[c] = new_cost
					came_from[c] = current.pos

					if debug {
						fmt.Printf("New node to visit: %v\n", wc)
					}
				} else {
					if debug {
						fmt.Printf("New cost %v is higher than current cost %v\n", new_cost, cost_map[c])
					}
				}
			} else {
				if debug {
					fmt.Printf("Not valid because of grid or reversing direction.\n")
				}
			}
		}
	}
	return min(cost_map[Vector{target, HORIZONTAL}], cost_map[Vector{target, VERTICAL}])
	//return 0
}

func processWorker(id int, work <-chan Position, resultChan chan<- int, wg *sync.WaitGroup,
	grid map[Position]int, width, height int, target Position) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		//result := processMirrorImage(w, false)
		result := 0
		resultChan <- walkGrid(grid, width, height, w, target)
		fmt.Printf("Worker=%v: %v => %v\n", id, w, result)
	}
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
	grid := make(map[Position]int)
	height := 0
	width := 0
	for y, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			for x, r := range cleaned_line {
				position := Position{x, y}
				grid[position] = toInt(string(r))
				if x > width {
					width = x
				}
			}
			if y > height {
				height = y
			}
			if debug {
				fmt.Printf("%v %v \n", cleaned_line, len(cleaned_line))
			}
		}
	}
	height++
	width++
	//**************************+
	channel_length := 1
	//channel_length := height*2 + width*2
	resultChan := make(chan int, channel_length)
	workChan := make(chan Position, channel_length)
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
	target := Position{x: width - 1, y: height - 1}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processWorker(i, workChan, resultChan, &wg, grid, width, height, target)
	}
	fmt.Printf("To work on: %v starts with %v threads\n", channel_length, numWorkers)

	b := Position{0, 0}
	workChan <- b
	close(workChan)

	// Close the channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := 0
	for result := range resultChan {
		if result > results {
			results = result
		}
	}
	fmt.Println("Final results:", results)

	fmt.Println("")
}
