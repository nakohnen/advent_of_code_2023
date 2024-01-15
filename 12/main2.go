package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	//"unicode"
	"math"
	"math/big"

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

func sumIntSlice(slice []int) int {
	res := 0
	for _, v := range slice {
		res += v
	}
	return res
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func getUniqueIntSlices(input [][]int) [][]int {
	result := [][]int{}
	for _, slice := range input {
		found := false
		for _, r_slice := range result {
			if equal(slice, r_slice) {
				found = true
			}
		}
		if !found {
			result = append(result, slice)
		}
	}
	return result
}

func getUniqueIntValues(input []int) []int {
	seen := make(map[int]bool)
	res := []int{}
	for _, i := range input {
		if !seen[i] {
			res = append(res, i)
			seen[i] = true
		}
	}
	return res
}

func factorial(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

func binomial(n, k int64) *big.Int {
	// Since C(n, k) = n! / (k! * (n - k)!)
	result := big.NewInt(0)
	nFactorial := factorial(n)
	kFactorial := factorial(k)
	nkFactorial := factorial(n - k)

	// Calculate the denominator (k! * (n - k)!)
	denominator := big.NewInt(0)
	denominator.Mul(kFactorial, nkFactorial)

	// Calculate the binomial coefficient
	result.Div(nFactorial, denominator)

	return result
}

// Convert BigInt to int64 and check for overflow
func convertBigIntToInt64(n *big.Int) int64 {
	minInt64 := big.NewInt(math.MinInt64)
	maxInt64 := big.NewInt(math.MaxInt64)

	// Check if n is within the range of int64
	if n.Cmp(minInt64) >= 0 && n.Cmp(maxInt64) <= 0 {
		return n.Int64() // No overflow, safe to convert
	}
	panic("Overflow happened")
	//return 0, false // Overflow occurred
}

// ConvertInt64ToInt converts an int64 to an int and checks for overflow
func ConvertInt64ToInt(n int64) int {
	/*if n > math.MaxInt32 || n < math.MinInt32 {
	    // Check for 32-bit architecture
	    panic("Overflow occurred")
	    //return 0, false // Overflow occurred
	}*/
	return int(n) // Safe to convert
}

func distributeStonesInBags(stones, bags int) int {
	if stones < 0 || bags < 0 {
		return 0
	}
	result := binomial(int64(stones+bags-1), int64(bags-1))
	result64 := convertBigIntToInt64(result)
	return ConvertInt64ToInt(result64)
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
	} else {
		// Lets check if the first filled blocks are correct
		filled := []string{}
		for _, s := range spring_split {
			if len(s) > 0 {
				if countChar(s, '?') == 0 {
					filled = append(filled, s)
				} else {
					break
				}
			}
		}
		for i := 0; i < min(len(filled), len(config)); i++ {
			if config[i] != len(filled[i]) {
				if debug {
					fmt.Printf("%v %v is invalid. ", filled, config)
				}
				return false
			}
		}
	}

	return true
}

func calculatePossibilities(spring string, config []int) int {
	cases := [2]rune{'.', '#'}
	work := []string{spring}
	if debug {
		fmt.Printf("%v with %v\n", spring, config)
	}
	if countChar(spring, '?') == 0 && len(config) == 1 && config[0] == len(spring) {
		return 1
	}

	result := 0
	if countChar(spring, '#') == 0 && countChar(spring, '.') == 0 {
		bags := len(config) + 1
		stones := len(spring) - (sumIntSlice(config) + len(config) - 1)
		result = distributeStonesInBags(stones, bags)
		if debug {
			fmt.Printf("\tv4 spring=%v config=%v combinations=%v and bags=%v and stones=%v\n", spring, config, result, bags, stones)
		}
		return result
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
						if countChar(new_spring, '?') == 0 {
							result++
						} else {
							new_work = append(new_work, new_spring)
						}
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
	if debug {
		fmt.Println("Possibilities:")
		fmt.Printf("\t   %v\n", spring)
		for _, w := range work {
			fmt.Printf("\t=> %v\n", w)
		}
	}
	if debug {
		fmt.Printf("%v %v => %v\n", spring, config, result)
	}
	return result
}

func createPossibleConfigsFromConfig(spring string, config []int) [][]int {
	result := [][]int{}
	result = append(result, []int{})
	len_spring := len(spring)
	for length := 1; length <= len(config); length++ {
		for i := 0; i <= len(config)-length; i++ {
			candidate := []int{}
			candidate = append(candidate, config[i:i+length]...)
			if sumIntSlice(candidate) <= len_spring {
				result = append(result, candidate)
			}
		}
	}
	return getUniqueIntSlices(result)
}

func createPossibleConfigs(spring string, filter []int) [][]int {
	all_configs := [][]int{}
	new_filter := getUniqueIntValues(filter)
	new_filter = append(new_filter, 0)
	for _, f := range new_filter {
		cand := []int{f}
		all_configs = append(all_configs, cand)
	}
	for length := 1; length < len(spring); length++ {
		new_configs := [][]int{}
		for _, old := range all_configs {
			for _, f := range new_filter {
				cand := []int{}
				cand = append(cand, old...)
				cand = append(cand, f)
				if sumIntSlice(cand) <= len(spring) {
					new_configs = append(new_configs, cand)
				}
			}
		}
		all_configs = new_configs
	}
	new_configs := [][]int{}
	for _, c := range all_configs {
		cand := []int{}
		for _, val := range c {
			if val > 0 {
				cand = append(cand, val)
			}
		}
		new_configs = append(new_configs, cand)
	}
	return getUniqueIntSlices(new_configs)
}

func multiplyString(s string, times int, join_char rune) string {
	var builder strings.Builder
	for i := 0; i < times; i++ {
		builder.WriteString(s)
		if join_char != 0 && i < times-1 {
			builder.WriteRune(join_char)
		}
	}
	return builder.String()
}

func multiplySlice[T any](slice []T, times int) []T {
	new_s := []T{}
	for i := 0; i < times; i++ {
		new_s = append(new_s, slice...)
	}
	return new_s
}

func flattenSliceIntSlice(input [][]int) []int {
	result := []int{}
	for _, inner := range input {
		result = append(result, inner...)
	}
	return result
}

func isValidSubSuperConfig(super [][]int, config []int) bool {
	flat_super := flattenSliceIntSlice(super)
	if len(flat_super) > len(config) {
		return false
	}
	for i := 0; i < min(len(flat_super), len(config)); i++ {
		if flat_super[i] != config[i] {
			return false
		}
	}
	return true
}

func recreateSuperConfig(input []int, possible_configs [][][]int) [][]int {
	result := [][]int{}
	for i, val := range input {
		result = append(result, possible_configs[i][val])
	}
	return result
}

func isValidSuperConfig(super [][]int, config []int, splits []string) bool {
	flat_super := flattenSliceIntSlice(super)
	if len(config) != len(flat_super) || len(super) != len(splits) {
		return false
	}
	for i := 0; i < len(config); i++ {
		if flat_super[i] != config[i] {
			return false
		}
	}
	return true
}

func findRunePosition(str string, r rune) int {
	for i, v := range str {
		if v == r {
			return i // Return the index if the rune is found
		}
	}
	return -1 // Return -1 if the rune is not found
}

func isSubSlice(sub, super []int) bool {
	return false
}

func createAllIntSliceDistributions(bins, stones int, max_bins []int) [][]int {
	result := [][]int{}
	for i := 0; i <= max_bins[0]; i++ {
		inner := []int{i}
		result = append(result, inner)
	}
	fmt.Printf("round=%v => %v\n", 0, result)
	for i := 1; i < bins; i++ {
		new_result := [][]int{}
		for _, inner := range result {
			todo := max(0, min(stones-sumIntSlice(inner), max_bins[i]))
			if debug {
				fmt.Printf("We have slice %v and we should add %v\n", inner, todo)
			}
			new_inner := []int{}
			for _, i_v := range inner {
				new_inner = append(new_inner, i_v)
			}
			for j := 0; j <= todo; j++ {
				new_inner_copy := []int{}
				for _, i_v := range new_inner {
					new_inner_copy = append(new_inner_copy, i_v)
				}
				new_inner_copy = append(new_inner_copy, j)
				new_result = append(new_result, new_inner_copy)

				new_inner_total := sumIntSlice(new_inner_copy)
				if new_inner_total > stones {
					if debug {
						fmt.Printf("bins=%v stones=%v new_inner_total=%v\n", bins, stones, new_inner_total)
						fmt.Printf("%v\n", new_inner_copy)
					}
					panic("We broke an assumption about bin distributions.")
				}
			}
		}
		result = new_result
		if debug {
			if len(result) < 100 {
				fmt.Printf("round=%v => %v\n", i, result)
			} else {
				fmt.Printf("round=%v => %v\n", i, len(result))
			}
		}
	}
	if debug {
		fmt.Println("We calculated all possible distributions")
	}
	return result
}

func createConfigDistributions(splits []string, config []int) [][][]int {
	result := [][][]int{}
	max_bins := []int{}
	for _, s := range splits {
		max_bins = append(max_bins, len(s))
	}
	for _, int_distribution := range createAllIntSliceDistributions(len(splits), len(config), max_bins) {
		new_config_d := [][]int{}
		todo := []int{}
		for _, c := range config {
			todo = append(todo, c)
		}
		for _, count := range int_distribution {
			inner := []int{}
			for i := 0; i < count; i++ {
				inner = append(inner, todo[0])
				todo = todo[1:]
			}
			new_config_d = append(new_config_d, inner)
		}
		valid := true
		for i, config_part := range new_config_d {
			total_val := sumIntSlice(config_part)
			if total_val > len(splits[i]) {
				valid = false
			} else if total_val == 0 && countChar(splits[i], '#') > 0 {
				valid = false
			}
		}
		if valid {
			result = append(result, new_config_d)
		}

	}
	return result
}

func calculateSplitsPossibilities(splits []string, config []int, memo *MemoizedMap) int {
	// Calculate all possible configs
	//   first level is the split => for each split a [][]int{}
	//   second level is the possible configs which have a []int{} shape
	possible_configs := [][][]int{}
	possible_combinations := [][]int{}
	for _, split := range splits {
		new_configs_tmp := createPossibleConfigsFromConfig(split, config)
		//new_configs_tmp := createPossibleConfigs(split, config)
		combinations := []int{}
		new_configs := [][]int{}
		impossible := [][]int{}
		for i := 0; i < len(new_configs_tmp); i++ {
			new_config := new_configs_tmp[i]
			pos := calculateTotalPossibilities(split, new_config, memo)
			if pos > 0 {
				// Only add those with make sense and are possible
				combinations = append(combinations, pos)
				new_configs = append(new_configs, new_config)
			} else {
				impossible = append(impossible, new_config)
			}
		}

		possible_configs = append(possible_configs, new_configs)
		possible_combinations = append(possible_combinations, combinations)
	}
	if debug {

		fmt.Printf("%v with config %v =>\n", splits, config)
		for i := 0; i < len(possible_configs); i++ {
			fmt.Printf("\t%v => %v == %v\n", splits[i], possible_configs[i],
				possible_combinations[i])
		}
		//fmt.Println("End of debug message.")
	}

	// Combine configs to calculate all possibilities and return
	super_configs := [][]int{}
	//fmt.Println("Here 0")
	for i := range possible_configs[0] {
		new_super := []int{i}
		super_configs = append(super_configs, new_super)
	}
	//fmt.Println("Here 1")
	result := 0
	//fmt.Println("Step 0")
	for i := 1; i < len(splits); i++ {
		new_outer_super_configs := [][]int{}
		for j := range possible_configs[i] {
			for _, super := range super_configs {
				new_super := []int{}
				new_super = append(new_super, super...)
				new_super = append(new_super, j)

				//fmt.Println("Step 1")
				super_with_config := recreateSuperConfig(new_super, possible_configs)
				if isValidSuperConfig(super_with_config, config, splits) {
					combs := 1

					//fmt.Println("Step 2")
					for split_i, inner_index := range new_super {
						combs *= possible_combinations[split_i][inner_index]
					}
					if debug {
						fmt.Printf("\tsuper=%v for config=%v == %v combinations\n", super_with_config, config, combs)
					}
					result += combs
				} else if isValidSubSuperConfig(super_with_config, config) {
					new_outer_super_configs = append(new_outer_super_configs, new_super)
				}
			}
			/*if debug {
			    fmt.Printf("\tsuper_configs=%v\n", super_configs)
			}*/
		}
		super_configs = new_outer_super_configs
	}
	/*if debug {
	    for _, super := range super_configs {
	        fmt.Printf("\tsuper=%v\n", super)
	    }
	}*/
	if debug {
		fmt.Printf("\tv2 spring=%v config=%v combinations=%v\n", splits, config, result)
	}
	return result
}

func calculateSimpleSplitsPossibilities(splits []string, config []int, memo *MemoizedMap) int {
	config_distributions := createConfigDistributions(splits, config)
	result := 0
	for _, config_candidate := range config_distributions {
		inner := 1
		for i, values := range config_candidate {
			pos := calculateTotalPossibilities(splits[i], values, memo)
			inner *= pos
			if pos == 0 {
				inner = 0
				break
			}
		}
		result += inner

	}
	return result
}

// intSliceToString converts an int slice to a string representation
func intSliceToString(slice []int) string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, v := range slice {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("]")

	return sb.String()
}

type MemoizedMap struct {
	calls    int
	mut      sync.Mutex
	memo     map[string]int
	memo_set map[string]bool
}

func calculateTotalPossibilities(spring string, config []int, memo *MemoizedMap) int {
	if sumIntSlice(config) > len(spring) {
		return 0
	}
	if len(config) == 0 {
		if countChar(spring, '#') > 0 {
			return 0
		} else {
			return 1
		}
	}
	if len(spring) < 4 || (countChar(spring, '?') < 4 && countChar(spring, '#') > 0) {
		return calculatePossibilities(spring, config)
	}
	// We split it along '.'
	splits_tmp := strings.Split(spring, ".")
	splits := []string{}
	for _, s := range splits_tmp {
		if len(s) > 0 {
			splits = append(splits, s)
		}
	}
	if len(splits) == 1 {
		if countChar(splits[0], '#') == 0 {
			new_spring := splits[0]
			bags := len(config) + 1
			stones := len(new_spring) - (sumIntSlice(config) + len(config) - 1)
			if debug {
				fmt.Printf("%v len=%v, %v bags, %v stones\n", new_spring, len(new_spring), bags, stones)
			}
			result := distributeStonesInBags(stones, bags)
			if debug {
				fmt.Printf("\tv3 spring=%v config=%v combinations=%v\n", new_spring, config, result)
			}
			return result

		} else {
			// Replace first '#' by '?' - '.'
			new_spring := splits[0]
			index := findRunePosition(new_spring, '#')
			//fmt.Printf("Index %v of # in %v\n", index, new_spring)
			full_spring := new_spring[:index] + "?"
			complement_spring := new_spring[:index] + "."
			if index+1 < len(new_spring) {
				full_spring += new_spring[index+1:]
				complement_spring += new_spring[index+1:]
			}
			if debug {
				fmt.Printf("Splitting %v into %v and %v\n", new_spring, full_spring, complement_spring)
			}
			result := memoizedCalculateTotPos(full_spring, config, memo) - memoizedCalculateTotPos(complement_spring, config, memo)
			//result := calculateTotalPossibilities(full_spring, config) - calculateTotalPossibilities(complement_spring, config)
			if debug {
				fmt.Printf("\tv5 spring=%v config=%v combinations=%v\n", spring, config, result)
			}
			return result
		}
	}

	if countChar(spring, '?') < 4 {
		return calculatePossibilities(spring, config)
	}
	//return calculateSimpleSplitsPossibilities(splits, config, memo)
	return calculateSplitsPossibilities(splits, config, memo)
}

func memoizedCalculateTotPos(spring string, config []int, memo *MemoizedMap) int {
	result := 0
	key := spring + intSliceToString(config)

	// Calculate symmetric position
	/*var sb strings.Builder
	for i := len(spring) - 1; i >= 0; i-- {
		sb.WriteString(string(spring[i]))
	}
	reverse_spring := sb.String()
	reverse_config := []int{}
	for i := len(config) - 1; i >= 0; i-- {
		reverse_config = append(reverse_config, config[i])
	}
	reverse_key := reverse_spring + intSliceToString(reverse_config) */

	// Check if we find something memoized
	found := false
	memo.mut.Lock()
	memo.calls += 1
	//if memo.memo_set[key] || memo.memo_set[reverse_key] {
	if memo.memo_set[key] { 
		result = memo.memo[key]
        if debug {
		    fmt.Printf("Found cached %v = %v (total calls = %v)\n", key, result, memo.calls)
        }
        found = true
	}
	memo.mut.Unlock()

	// If not found calculate it and save it
	if !found {
		result = calculateTotalPossibilities(spring, config, memo)
		memo.mut.Lock()
		memo.memo_set[key] = true
		//memo.memo_set[reverse_key] = true
		memo.memo[key] = result
		//memo.memo[reverse_key] = result
		memo.mut.Unlock()
	}
	return result
}

type WorkItem struct {
	spring string
	config []int
}

func processWorker(id int, work <-chan WorkItem, resultChan chan<- int, wg *sync.WaitGroup, memo *MemoizedMap) {
	defer wg.Done()
	for w := range work {
		//fmt.Printf("Worker=%v calculating %v %v\n", id, w.spring, w.config)
		pos := memoizedCalculateTotPos(w.spring, w.config, memo)
		resultChan <- pos
        fmt.Printf("Worker=%v: %v %v => %v\n", id, w.spring, w.config, pos)
	}
}

func main() {

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

	springs := []string{}
	configs := [][]int{}
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	for _, line := range lines {
		line_split := strings.Split(line, " ")
		spring := line_split[0]
		spring = multiplyString(spring, 5, '?')
		springs = append(springs, spring)

		config := []int{}
		config_part := line_split[1]
		config_part = multiplyString(config_part, 5, ',')
		for _, val := range strings.Split(config_part, ",") {
			config = append(config, toInt(val))
		}
		configs = append(configs, config)
		fmt.Printf("%v => %v + %v\n", line, spring, config)
	}

	results := 0
	var memo MemoizedMap = MemoizedMap{
		memo:     make(map[string]int),
		memo_set: make(map[string]bool),
	}
	if !debug {
		resultChan := make(chan int, len(springs))
		workChan := make(chan WorkItem, len(springs))
		var wg sync.WaitGroup

		var numWorkers int = *cores
		if numWorkers < 1 {
			numWorkers = 1
		} else if numWorkers > 24 {
			numWorkers = 24
		}
		numWorkers = min(numWorkers, len(springs))

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go processWorker(i, workChan, resultChan, &wg, &memo)
		}
		fmt.Printf("To work on: %v elements with %v threads\n", len(springs), numWorkers)

		for i := 0; i < len(springs); i++ {
			workChan <- WorkItem{spring: springs[i], config: configs[i]}

		}
		close(workChan)

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
			results += memoizedCalculateTotPos(springs[i], configs[i], &memo)
		}
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
