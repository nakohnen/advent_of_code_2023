package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	//"unicode"
	//"math"
	//"math/big"

	//"sort"
	"strconv"
	//"sync"
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

func transposeStringSlice(slice []string) []string {
	rows := len(slice)
	cols := len(slice[0])
	result := []string{}
	if debug {
		fmt.Printf("Rows=%v cols=%v len(slice)=%v\n", rows, cols, len(slice))
	}
	var sb strings.Builder
	for col := 0; col < cols; col++ {
		sb.Reset()
		for row := 0; row < rows; row++ {
			sb.WriteByte(slice[row][col])
		}
		result = append(result, sb.String())
	}
	if debug {
		for i := 0; i < max(len(slice), len(result)); i++ {
			if i < len(slice) {
				fmt.Printf("%v ", slice[i])
			} else {
				sb.Reset()
				for j := 0; j < len(slice[0]); j++ {
					sb.WriteString(" ")
				}
				fmt.Printf("%v ", sb.String())
			}
			if i < len(result) {
				fmt.Printf("%v\n", result[i])
			} else {
				fmt.Println("")
			}
		}
	}

	return result
}

func countChar(s string, c rune) int {
	r := 0
	for _, sc := range s {
		if sc == c {
			r++
		}
	}
	return r
}

func sumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
}

func rollRocksAndCalculate(rocks string) int {
    width := len(rocks)
	if debug {
		fmt.Printf("Roll rocks: %v\n with width %v", rocks, width)
	}
	sections := strings.Split(rocks, "#")
	var sb strings.Builder
	new_sections := []string{}
	for _, section := range sections {
		sb.Reset()
		len_section := len(section)
		stones := countChar(section, 'O')
		for i := 0; i < len_section; i++ {
			r := '.'
			if i < stones {
				r = 'O'
			}
			sb.WriteRune(r)
		}
		new_sections = append(new_sections, sb.String())
	}
	new_rocks := strings.Join(new_sections, "#")
	if debug {
		fmt.Printf("Calculate load: %v\n", new_rocks)
	}
	result := 0
	for i, r := range new_rocks {
		if r == 'O' {
			result += width - i
		}

	}
    fmt.Println(new_rocks)

	return result
}

func processWorker(id int, work <-chan string, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		//result := processMirrorImage(w, false)
		fmt.Printf("\n")
		result := rollRocksAndCalculate(w)
		resultChan <- result
		fmt.Printf("Worker=%v: %v => %v\n", id, w, result)
	}
}

func main() {
	if debug {
		fmt.Println("Debug enabled")
	}

	cores := flag.Int("t", 1, "On how many concurrent goroutines should this code run? (1-24)")
	filename_f := flag.String("f", "", "On which file should we run this?")
	flag.Parse()

	// The second element in os.Args is the first argument
	filename := *filename_f
	if len(filename) == 0 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

	dat, err := os.ReadFile(filename)
	check(err)

	// Read file
	image := []string{}
	width := 0
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			width = len(cleaned_line)
			image = append(image, cleaned_line)
			fmt.Printf("%v\n", cleaned_line)
		}
	}
	fmt.Printf("Transpose image: width=%v height=%v\n", width, len(image))
	image = transposeStringSlice(image)
	//**************************+

	resultChan := make(chan int, len(image))
	workChan := make(chan string, len(image))
	var wg sync.WaitGroup

	var numWorkers int = *cores
	if debug {
		numWorkers = 1
	} else {
		if numWorkers < 1 {
			numWorkers = 1
		} else if numWorkers > 24 {
			numWorkers = 24
		}
		numWorkers = min(numWorkers, len(image))
	}

	fmt.Println("Work on image.")
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processWorker(i, workChan, resultChan, &wg)
	}
	fmt.Printf("To work on: %v mirrors with %v threads\n", len(image), numWorkers)

	for _, m := range image {
		workChan <- m

	}
	close(workChan)

	// Close the channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := 0
	for result := range resultChan {
		results += result
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
