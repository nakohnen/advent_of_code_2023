package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"unicode"
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

func sumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
}

func hashString(s string) int {
	result := 0
	for _, r := range s {
		if unicode.IsPrint(r) {
			fmt.Printf("%v = %v\n", r, string(r))
			result += int(r)
			result *= 17
			result %= 256
		}
	}

	return result
}

func processWorker(id int, work <-chan string, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		//result := processMirrorImage(w, false)
		fmt.Printf("\n")
		result := hashString(w)
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
	instructions := []string{}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)

		for _, sub := range strings.Split(cleaned_line, ",") {
			if len(sub) > 0 {
				instructions = append(instructions, sub)
			}
		}
	}
	//**************************+

	resultChan := make(chan int, len(instructions))
	workChan := make(chan string, len(instructions))
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
		numWorkers = min(numWorkers, len(instructions))
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processWorker(i, workChan, resultChan, &wg)
	}
	fmt.Printf("To work on: %v instructions with %v threads\n", len(instructions), numWorkers)

	for _, m := range instructions {
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
