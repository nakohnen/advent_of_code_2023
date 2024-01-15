package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	//"unicode"
	//	"math"
	//"sort"
	"strconv"
	//"sync"
)

const debug bool = true

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

func countConseqChar(spring string, char rune) int {
	start := false
	res := 0
	for _, r := range spring {
		if r == char {
			if !start {
				start = true
			}
			res++
		} else if r != char {
			if start {
				return res
			}
		}
	}
	return res

}

func countChar(spring string, char rune) int {
	res := 0
	for _, r := range spring {
		if r == char {
			res++
		}
	}
	return res
}

func checkCorrectness(spring string, config []int) bool {
	sum_c := 0
	for _, v := range config {
		sum_c += v
	}
	if countChar(spring, '#')+countChar(spring, '?') < sum_c {
		return false
	}
	if countChar(spring, '#') > sum_c {
		return false
	}

	no_q_marks := true
	for _, r := range spring {
		if r == '?' {
			no_q_marks = false
			break
		}
	}

	spring_split := strings.Split(spring, ".")
	if no_q_marks {
		p_config := []int{}
		for _, s := range spring_split {
			if len(s) > 0 {
				p_config = append(p_config, len(s))
			}
		}
		if len(p_config) == len(config) {
			for i := 0; i < len(config); i++ {
				if p_config[i] != config[i] {
					return false
				}
			}
			return true
		}
		return false
	}

	return true
}

func sumIntSlice(slice []int) int {
	res := 0
	for _, v := range slice {
		res += v
	}
	return res
}

func calculatePossibilites(spring string, config []int) int {
	cases := [2]rune{'.', '#'}
	work := []string{spring}
	if debug {
		fmt.Printf("%v with %v\n", spring, config)
	}
	for i := 0; i < len(spring); i++ {
		new_work := []string{}
		for _, w := range work {
			if w[i] == '?' {
				for _, r := range cases {
					runes := []rune(w)
					runes[i] = r
					new_spring := string(runes)
					if debug {
						fmt.Printf("New candidate %v", new_spring)
					}

					if countChar(new_spring, '?') > 0 && countChar(new_spring, '#') == sumIntSlice(config) {
						new_spring2 := replaceCharacters(new_spring, "?", ".")
						if debug {
							fmt.Printf(" => Replacing %v with %v", new_spring, new_spring2)
						}
						new_spring = new_spring2
					}

					if checkCorrectness(new_spring, config) {
						if debug {
							fmt.Print(" is valid.")
						}
						new_work = append(new_work, new_spring)
					}
					if debug {
						fmt.Println("")
					}
				}
			} else {
				new_work = append(new_work, w)
			}
		}
		work = new_work
	}
	result := len(work)
	if debug {
		fmt.Println("Possibilites:")
		fmt.Printf("\t   %v\n", spring)
		for _, w := range work {
			fmt.Printf("\t=> %v\n", w)
		}
	}
	fmt.Printf("%v %v => %v\n", spring, config, result)
	return result
}

func processWorker(spring string, config []int, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	resultChan <- calculatePossibilites(spring, config)
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

	springs := []string{}
	configs := [][]int{}
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	for _, line := range lines {
		line_split := strings.Split(line, " ")
		springs = append(springs, line_split[0])

		config := []int{}
		for _, val := range strings.Split(line_split[1], ",") {
			config = append(config, toInt(val))
		}
		configs = append(configs, config)
		fmt.Printf("%v => %v + %v\n", line, line_split[0], config)
	}

	results := 0
	if !debug {
		resultChan := make(chan int, len(springs))
		var wg sync.WaitGroup

		for i := 0; i < len(springs); i++ {
			wg.Add(1)
			go processWorker(springs[i], configs[i], resultChan, &wg)
		}
		fmt.Printf("To work on: %v elements\n", len(springs))

		// Close the channel once all goroutines have finished
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// Collect results
		for result := range resultChan {
			results += result
		}
	} else {
		for i := 0; i < len(springs); i++ {
			results += calculatePossibilites(springs[i], configs[i])
		}
	}

	fmt.Println("Results:", results)

	fmt.Println("")
}
