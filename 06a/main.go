package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
	//	"math"
	"strconv"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hasPrefix(line string, sub string) bool {
	if len(line) < len(sub) {
		return false
	}
	return line[0:len(sub)] == sub
}

func main() {
	//dat, err := os.ReadFile("input_sample.txt")
	dat, err := os.ReadFile("input.txt")
	check(err)

	time := []int{}
	distance := []int{}

	for _, line := range strings.Split(strings.TrimSpace(string(dat)), "\n") {
		fmt.Printf("%v => ", line)

		var working_slice []int

		line_split := strings.Split(line, " ")
		for _, n := range line_split[1:] {
			if n != "" {
				v, err := strconv.Atoi(n)
				check(err)
				if v > 0 {
					working_slice = append(working_slice, v)
				}
			}
		}
		fmt.Printf("%v \n", working_slice)
		if hasPrefix(line, "Time:") {
			time = working_slice
		} else if hasPrefix(line, "Distance:") {
			distance = working_slice
		}

	}
    // reconstruct time

	fmt.Printf("Time: %v\n", time)
	fmt.Printf("Distance: %v\n", distance)
	acc := 1
	possibilites := []int{}
	for i := 0; i < len(time); i++ {
		fmt.Printf("%v %v %v => ", i, time[i], distance[i])
		n := time[i] - 1
		sum := 0
		var wg sync.WaitGroup
		wg.Add(n)

		results := make(chan int, n)

		task := func(t int) {
			run_time := time[i] - t
			result := 0
			if run_time*t > distance[i] {
				result = 1
			}
			results <- result
			wg.Done()
		}

		for j := 0; j < n; j++ {
			go task(j)
		}

		wg.Wait()
		close(results)
		for r := range results {
			sum += r
		}
		if sum > 0 {
			acc *= sum
		}
		possibilites = append(possibilites, sum)
		fmt.Printf("%v -> %v\n", sum, acc)
	}

	fmt.Println("")
}
