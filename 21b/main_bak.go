package main

import (
	"flag"
	"fmt"
	"math"
	"os"
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

func allElementsSame[T comparable](slice []T) bool {
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			return false
		}
	}
	return true
}

// Function to find prime factors of a number
func primeFactors(n int) map[int]int {
	factors := make(map[int]int)
	// Count the number of 2s that divide n
	for n%2 == 0 {
		factors[2]++
		n = n / 2
	}
	// n must be odd at this point. So start from 3 and iterate until sqrt(n)
	for i := 3; i <= int(math.Sqrt(float64(n))); i = i + 2 {
		// While i divides n, count i and divide n
		for n%i == 0 {
			factors[i]++
			n = n / i
		}
	}
	// If n is a prime number greater than 2
	if n > 2 {
		factors[n]++
	}
	return factors
}

// Function to find LCM of an array of integers
func findLCM(arr []int) int {
	overallFactors := make(map[int]int)
	for _, num := range arr {
		// Get prime factors of each number
		primeFactorsOfNum := primeFactors(num)
		for prime, power := range primeFactorsOfNum {
			if currentPower, exists := overallFactors[prime]; !exists || power > currentPower {
				// Store the highest power of each prime
				overallFactors[prime] = power
			}
		}
	}
	// Calculate LCM by multiplying the highest powers of all primes
	lcm := 1
	for prime, power := range overallFactors {
		lcm *= int(math.Pow(float64(prime), float64(power)))
	}
	return lcm
}

func leastCommonMultiple(slice []int) int {
	calc := []int{}
	for _, i := range slice {
		calc = append(calc, i)
	}
	for !allElementsSame[int](calc) {
		min_element := math.MaxInt
		for _, v := range calc {
			min_element = min(v, min_element)
		}
		min_index := IndexOf[int](calc, min_element)
		calc[min_index] += slice[min_index]
		fmt.Printf("LCM (%v): %v\n", slice, calc)
	}

	return calc[0]
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

func removeDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, i := range slice {
		if !seen[i] {
			result = append(result, i)
			seen[i] = true
		}
	}
	return result
}

type Position struct {
	x, y int
}

func addPositions(a, b Position) Position {
	return Position{x: a.x + b.x, y: a.y + b.y}
}

type PlotType bool

const (
	GARDEN PlotType = false
	ROCK   PlotType = true
)

func getAdjacent(current, start Position, width, heigth int, plots_map map[Position]PlotType, steps int) []Position {
	result := []Position{}
	candidates := []Position{}
	candidates = append(candidates, Position{x: current.x - 1, y: current.y})
	candidates = append(candidates, Position{x: current.x + 1, y: current.y})
	candidates = append(candidates, Position{x: current.x, y: current.y - 1})
	candidates = append(candidates, Position{x: current.x, y: current.y + 1})
	for _, c := range candidates {
		if c.x >= -steps-start.x && c.x <= steps+start.x && c.y >= -steps-start.y && c.y <= steps+start.y {
			new_x := c.x
			new_y := c.y
			if c.x < 0 {
				new_x = width - abs(c.x)%width
			} else if c.x >= width {
				new_x = c.x % width
			}
			if c.y < 0 {
				new_y = heigth - abs(c.y)%heigth
			} else if c.y >= heigth {
				new_y = c.y % heigth
			}
			if !plots_map[Position{x: new_x, y: new_y}] {
				if debug {
					fmt.Printf("%v: translated %v, %v \n", c, new_x, new_y)
				}
				result = append(result, c)
			}
		}
	}
	return result
}

func getAdjacentConstrained(current Position, width, heigth int, plots_map map[Position]PlotType) []Position {
	result := []Position{}
	candidates := []Position{}
	candidates = append(candidates, Position{x: current.x - 1, y: current.y})
	candidates = append(candidates, Position{x: current.x + 1, y: current.y})
	candidates = append(candidates, Position{x: current.x, y: current.y - 1})
	candidates = append(candidates, Position{x: current.x, y: current.y + 1})
	for _, c := range candidates {
		if !plots_map[c] && c.x >= 0 && c.x < width && c.y >= 0 && c.y < heigth {
			result = append(result, c)
		}
	}
	return result
}
func getMinSteps(start Position, width, heigth int, plots_map map[Position]PlotType, steps int) map[Position]int {
	result := make(map[Position]int)
	for x := -steps - start.x; x <= steps+start.x; x++ {
		for y := -steps - start.y; y <= steps+start.y; y++ {
			result[Position{x: x, y: y}] = math.MaxInt
		}
	}
	work := []Position{start}
	min_steps := 0
	for len(work) > 0 {
		min_steps++
		new_work := []Position{}
		for _, current := range work {
			candidates := getAdjacent(current, start, width, heigth, plots_map, steps)
			for _, cand := range candidates {
				if min_steps < result[cand] && min_steps <= steps {
					new_work = append(new_work, cand)
					result[cand] = min_steps
				}
			}
		}
		work = removeDuplicates[Position](new_work)
	}
	return result
}

func getMinStepsConstrained(start Position, width, heigth int, plots_map map[Position]PlotType) map[Position]int {
	result := make(map[Position]int)
	for x := 0; x < width; x++ {
		for y := 0; y < width; y++ {
			result[Position{x, y}] = math.MaxInt
		}
	}
	work := []Position{start}
	steps := 0
	for len(work) > 0 {
		steps++
		new_work := []Position{}
		for _, current := range work {
			candidates := getAdjacentConstrained(current, width, heigth, plots_map)
			for _, cand := range candidates {
				if steps < result[cand] {
					new_work = append(new_work, cand)
					result[cand] = steps
				}
			}
		}
		work = removeDuplicates[Position](new_work)
	}
	return result
}

func checkerDistance(a, b Position) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func calculateReachableTiles(start, target_frame Position, width, height int, min_steps_frames map[Position]map[Position]int, frame_points []Position, max_steps int) int {
	min_distance := math.MaxInt
	min_point := target_frame
	for _, fp := range frame_points {
		new_fp := addPositions(target_frame, fp)
		d := checkerDistance(start, new_fp)
		if d < min_distance {
			min_distance = d
			min_point = fp
		}
	}
	//  fmt.Printf("Frame %v with distance %v and frame point %v\n", target_frame, min_distance, min_point)
	result := 0
	if min_distance > max_steps {
		//      fmt.Printf("Out!\n")
		return result
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pos := Position{x: x, y: y}
			t_steps := min_steps_frames[min_point][pos] + min_distance
			if t_steps <= max_steps && t_steps%2 == max_steps%2 {
				result++
			}
		}
	}
	//    fmt.Printf("In: %v\n", result)

	return result
}

func processFrame(id int, start Position, target <-chan Position, width, height int, min_steps_larger map[Position]map[Position]int, frame_points []Position, max_steps int, result_chan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for t_frame := range target {
		// fmt.Printf("Frame %v (chan len=%v) \n", t_frame, len(target))
		result := calculateReachableTiles(start, t_frame, width, height, min_steps_larger, frame_points, max_steps)
        fmt.Printf("frame: %v => %v (max_steps=%v)\n", t_frame, result, max_steps)
		result_chan <- result
	}
}

func resizeStringSlice(slice []string, resize int) []string {
    res := []string{}
    for i:=0;i<resize;i++{
        for _, s := range slice {
            res = append(res, s)
        }
    }
    return res
}

func resizeString(s string, resize int) string {
    var sb strings.Builder
    for i:=0;i<resize;i++ {
        sb.WriteString(s)
    }
    return sb.String()
}

func main() {
	if debug {
		fmt.Println("Debug enabled")
	}
	cores_f := flag.Int("t", -1, "How many cores (threads) should we run?")
	filename_f := flag.String("f", "", "On which file should we run this?")
	steps_f := flag.Int("s", 64, "How many steps do we need to simulate?")
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

	steps := *steps_f
	if steps < 0 {
		steps = 1
	}

	dat, err := os.ReadFile(filename)
	check(err)

	// Read file
	plots := [][]PlotType{}
	start := Position{-1, -1}
    plots_map := make(map[Position]PlotType)
    
    resize_factor := 10
    
    starts := []Position{}
	for y, line := range resizeStringSlice(strings.Split(string(dat), "\n"), resize_factor) {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			//fmt.Printf("%v \n", cleaned_line)
			plots_row := []PlotType{}
			for x, s := range resizeString(cleaned_line, resize_factor) {
				pos := Position{x: x, y: y}
				tile := GARDEN
				if s == '#' {
					tile = ROCK
				} else if s == 'S' && x == y {
                    starts = append(starts, pos)
				}
				plots_map[pos] = tile
				plots_row = append(plots_row, tile)
			}
			plots = append(plots, plots_row)
		}
	}
	width := len(plots[0])
	heigth := len(plots)
    min_distance := math.MaxInt
    mid_point := Position{x: width/2, y: heigth/2}
    for _, s_i := range starts {
        d := checkerDistance(s_i, mid_point)
        if d < min_distance {
            start = s_i
            min_distance = d
        }
    }

	fmt.Printf("Simulating %v steps with start %v.\n", steps, start)
	fmt.Printf("Width=%v, height=%v\n", width, heigth)
	min_steps_larger := make(map[Position]map[Position]int)
	x_positions := [3]int{0, start.x, width - 1}
	y_positions := [3]int{0, start.y, heigth - 1}
	frame_points := []Position{}
	for _, x := range x_positions {
		for _, y := range y_positions {
			pos := Position{x: x, y: y}
			fmt.Printf("Calculating frame point %v", pos)
			frame_points = append(frame_points, pos)
			min_steps_larger[pos] = getMinStepsConstrained(pos, width, heigth, plots_map)
			fmt.Println(" ->  done!")
		}
	}
	fmt.Println("Done calculating minimum steps.")

	var wg sync.WaitGroup
	channel_width := max(steps/width, steps/heigth)

	var numWorkers int = cores
	if debug {
		numWorkers = 1
	} else {
		if numWorkers < 1 {
			numWorkers = 1
		} else if numWorkers > 24 {
			numWorkers = 24
		}
		numWorkers = max(1, min(numWorkers, channel_width))
	}

	// Prepare channels
	frames := make(chan Position, channel_width)
	result_chan := make(chan int, numWorkers)
	sum_result_chan := make(chan int)

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processFrame(i, start, frames, width, heigth, min_steps_larger, frame_points, steps, result_chan, &wg)
	}
	fmt.Printf("To work on: %v frames with %v threads\n", channel_width, numWorkers)

	// Start a goroutine to sum the results
	go func() {
		sum := 0
		for num := range result_chan {
			sum += num
		}
		sum_result_chan <- sum
		close(sum_result_chan)
	}()

	// Close outputChannel once all goroutines are done
	go func() {
		wg.Wait()
		close(result_chan)
	}()

	// Fill input queue
	final_frames_count := 0
	for x := 0; x < 2+steps/width; x++ {
		for y := 0; y < 2+steps/heigth; y++ {
			x_min := x * width
			y_min := y * heigth
			x_range := []int{-1, 1}
			if x == 0 {
				x_range = []int{1}
			}
			for _, x_m := range x_range {
				y_range := []int{-1, 1}
				if y == 0 {
					y_range = []int{1}
				}
				for _, y_m := range y_range {
					final_frames_count++
					frames <- Position{x: x_min * x_m, y: y_min * y_m}
				}
			}
		}
	}
	fmt.Printf("Final frames count vs prediction: %v vs %v\n", final_frames_count, channel_width)
	close(frames)

	results := <-sum_result_chan
	fmt.Println("Final results:", results)

	fmt.Println("")
}
