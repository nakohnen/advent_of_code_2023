package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

func checkerDistance(a, b Position) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
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

func resizeStringSlice(slice []string, resize int) []string {
	res := []string{}
	for i := 0; i < resize; i++ {
		for _, s := range slice {
			res = append(res, s)
		}
	}
	return res
}

func resizeString(s string, resize int) string {
	var sb strings.Builder
	for i := 0; i < resize; i++ {
		sb.WriteString(s)
	}
	return sb.String()
}

type Position struct {
	x, y int
}

type PlotType bool

const (
	GARDEN PlotType = false
	ROCK   PlotType = true
)

func getAdjacent(current Position, width, heigth int, plots_map map[Position]PlotType) []Position {
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

func getMinSteps(start Position, width, heigth int, plots_map map[Position]PlotType) map[Position]int {
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
			candidates := getAdjacent(current, width, heigth, plots_map)
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
	plots_map := make(map[Position]PlotType)
    starts := []Position{}

    raw_map := []string{}
    for _, line := range strings.Split(string(dat), "\n") { 
		cleaned_line := strings.TrimSpace(line)
        if len(cleaned_line) > 0 {
            raw_map = append(raw_map, cleaned_line)
        }
    }
	resize_f := 1 + 2 * (steps / (len(raw_map) / 2))
    fmt.Printf("Size = %v^2 => resize factor %v for steps %v\n", len(raw_map), resize_f, steps)
	for y, line := range resizeStringSlice(raw_map, resize_f) {
			//fmt.Printf("%v \n", cleaned_line)
			plots_row := []PlotType{}
			for x, s := range resizeString(line, resize_f) {
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
	width := len(plots[0])
	heigth := len(plots)
    min_d := math.MaxInt
    mid_point := Position{x:width/2, y:heigth/2}
    start := Position{-1, -1}
    for _, s_i := range starts {
        d := checkerDistance(s_i, mid_point)
        if d < min_d {
            min_d = d
            start = s_i
        }
    }
    fmt.Printf("Start=%v, width=%v, heigth=%v, steps=%v\n", start, width, heigth, steps)

	fmt.Printf("Simulating %v steps with start %v.\n", steps, start)
	min_steps := getMinSteps(start, width, heigth, plots_map)
	results := 0
    var sb strings.Builder
    for y := 0; y < heigth; y++ {
        sb.Reset()
	    for x := 0; x < width; x++ {
			pos := Position{x: x, y: y}
			min_step := min_steps[pos]
            is_valid := false
			if min_step <= steps && min_step%2 == steps%2 {
				results++
                is_valid = true
				if debug {
					// fmt.Printf("%v: %v %v\n", pos, min_step, steps)
				}
			}
            if is_valid {
                if plots_map[pos] == ROCK {
                    sb.WriteRune('X')
                } else {
                    sb.WriteRune('O')
                }
            } else {
                switch plots_map[pos] {
                case GARDEN:
                    sb.WriteRune('.')
                case ROCK:
                    sb.WriteRune('#')
                    
                }
            }
		}
        fmt.Println(sb.String())
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
