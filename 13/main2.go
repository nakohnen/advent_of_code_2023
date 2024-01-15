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

func flipOnePosition(mirror []string, col, row int) []string {
	result := []string{}
	if debug {
		fmt.Printf("Flip at col=%v row=%v\n", col, row)
	}
	for i := 0; i < len(mirror); i++ {
		current_row := mirror[i]
		if i == row {
			new_symbol := ""
			switch string(current_row[col]) {
			case ".":
				new_symbol = "#"
			case "#":
				new_symbol = "."
			}
			new_row := ""

			if col == 0 {
				new_row = new_symbol + current_row[1:]
			} else if col == len(current_row)-1 {
				new_row = current_row[:len(current_row)-1] + new_symbol
			} else {
				new_row = current_row[:col] + new_symbol + current_row[col+1:]
			}

			if len(new_row) != len(current_row) {
				panic("We dont have the same lengths.")
			}
			result = append(result, new_row)

			if debug {
				fmt.Printf("%v => %v\n", current_row, new_row)
			}
		} else {
			result = append(result, current_row)
			if debug {
				fmt.Printf("%v => %v\n", current_row, current_row)
			}
		}

	}
	if debug {
		fmt.Println("")
	}
	return result
}

func transposeStringSlice(slice []string) []string {
	rows := len(slice)
	cols := len(slice[0])
	result := []string{}
	/*if debug {
		fmt.Printf("Rows=%v cols=%v len(slice)=%v\n", rows, cols, len(slice))
	}*/
	for col := 0; col < cols; col++ {
		var sb strings.Builder
		for row := 0; row < rows; row++ {
			sb.WriteByte(slice[row][col])
		}
		result = append(result, sb.String())
	}
	/*if debug {
		for i := 0; i < max(len(slice), len(result)); i++ {
			if i < len(slice) {
				fmt.Printf("%v ", slice[i])
			} else {
				var sb strings.Builder
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
	}*/

	return result
}

func processMirrorImage(mirror []string, transposed bool, ignore int) int {
	candidates := []int{}
	last_row := mirror[0]
	for i := 1; i < len(mirror); i++ {
		current_row := mirror[i]
		/*if debug {
			fmt.Printf("%v vs %v\n", current_row, last_row)
		}*/
		if current_row == last_row && ignore != i {
			/*if debug {
				fmt.Printf("We found a candidate %v\n", i)
			}*/
			candidates = append(candidates, i)
		}
		last_row = current_row
	}
	/*if debug {
		fmt.Printf("We have candidates: %v\n", candidates)
	}*/

	for _, row_c := range candidates {
		inner_valid := true
		mirror_border := min(row_c, len(mirror)-row_c)

		/*if debug {
			fmt.Printf("Candidate row=%v height of mirror=%v border=%v\n", row_c, len(mirror), mirror_border)
		}*/

		for i := 0; i < mirror_border; i++ {
			mirror_c := row_c - i - 1
			target_c := row_c + i
			if mirror_c < 0 || target_c >= len(mirror) {
				/*if debug {
					fmt.Printf("Invalid because of invalid positions\n")
				}*/
				inner_valid = false
				break
			}
			/*if debug {
				fmt.Printf("%v vs %v\n", mirror[target_c], mirror[mirror_c])
			}*/
			if mirror[target_c] != mirror[mirror_c] {
				/*if debug {
					fmt.Printf("Invalid because not same.\n")
				}*/
				inner_valid = false
				break
			}
		}
		if inner_valid {
			if debug {
				fmt.Printf("We found a valid reflection at %v (t=%v).\n", row_c, transposed)
			}
			if transposed {
				return row_c
			} else {
				return 100 * (row_c)
			}
		}
	}
	return 0
}

func sumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
}

func fixSmudge(mirror []string) int {
	baseline := processMirrorImage(mirror, false, -1)
	if baseline == 0 {
		baseline = processMirrorImage(transposeStringSlice(mirror), true, -1)
	}
	if baseline > 0 {
		fmt.Printf("Baseline = %v\n", baseline)
	}
	ignore_h := -1
    ignore_v := -1
	if baseline >= 100 {
		ignore_h = baseline / 100
	} else {
        ignore_v = baseline
    }
	result := []int{}
	contained := make(map[int]bool)
	for row := 0; row < len(mirror); row++ {
		for col := 0; col < len(mirror[0]); col++ {
			fmt.Println("-------------------------------")
			flipped := flipOnePosition(mirror, col, row)
			r := processMirrorImage(flipped, false, ignore_h)
			r2 := processMirrorImage(transposeStringSlice(flipped), true, ignore_v)
			fmt.Printf("Result for col=%v row=%v => %v or %v\n", col, row, r, r2)
			if r > 0 && r != baseline {
				if !contained[r] {
					fmt.Printf("Fixed smudge at col=%v row=%v => %v\n", col, row, r)
					result = append(result, r)
					contained[r] = true
				}
			}
			if r2 > 0 && r2 != baseline {
				if !contained[r2] {
					fmt.Printf("Fixed smudge at col=%v row=%v => %v\n", col, row, r2)
					result = append(result, r2)
					contained[r2] = true
				}
			}
		}
	}
	return sumIntSlice(result)
}

func processWorker(id int, work <-chan []string, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		//result := processMirrorImage(w, false)
		fmt.Printf("\n")
		result := fixSmudge(w)
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
	mirrors := [][]string{}
	current_mirror := []string{}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if cleaned_line == "" {
			mirrors = append(mirrors, current_mirror)
			current_mirror = []string{}
		} else {
			current_mirror = append(current_mirror, cleaned_line)
		}
		fmt.Printf("%v\n", cleaned_line)
	}
	//**************************+

	resultChan := make(chan int, len(mirrors))
	workChan := make(chan []string, len(mirrors))
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
		numWorkers = min(numWorkers, len(mirrors))
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processWorker(i, workChan, resultChan, &wg)
	}
	fmt.Printf("To work on: %v mirrors with %v threads\n", len(mirrors), numWorkers)

	for _, m := range mirrors {
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
