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
    for i:=1;i<len(slice);i++{
        if slice[i] != slice[i-1]{
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

type ModuleType int

const (
	UNTYPED ModuleType = iota
	FLIPFLOP
	CONJUNCTION
)

func isAllLow(states map[string]bool, names []string) bool {
	for _, name := range names {
		if states[name] {
			return false
		}
	}
	return true
}

type Pulse struct {
	source  string
	target  string
	highlow bool
}

func processFlipflop(pulse Pulse, forward map[string][]string, states map[string]bool) []Pulse {
	result := []Pulse{}
	name := pulse.target
	if pulse.highlow == false {
		ff_state := !states[name]
		for _, other := range forward[name] {
			result = append(result, Pulse{name, other, ff_state})
		}
		states[name] = ff_state
	}
	return result
}

func processConjunction(pulse Pulse, forward, backward map[string][]string, conj_states map[string]map[string]bool) []Pulse {
	result := []Pulse{}
	name := pulse.target
	source := pulse.source

	conj_states[name][source] = pulse.highlow
	all_high := true
	for _, b_name := range backward[name] {
		if conj_states[name][b_name] == false {
			all_high = false
			break
		}
	}

	for _, t_name := range forward[name] {
		result = append(result, Pulse{name, t_name, !all_high})
	}

	return result
}

func processSignal(presses int, names []string, forward map[string][]string,
	backward map[string][]string, types map[string]ModuleType,
	states map[string]bool,	conj_states map[string]map[string]bool) []Pulse {
    result := []Pulse{}
	pulse := Pulse{"button", "broadcaster", false}
	to_work := []Pulse{pulse}
	for len(to_work) > 0 {
		current := to_work[0]
		to_work = to_work[1:]
		name := current.target

        if types[current.target] == CONJUNCTION && current.highlow == true {
            result = append(result, current)
        }
        others := []Pulse{}
		switch t := types[name]; t {
		case UNTYPED:
            for _, o := range forward[name] {
                others = append(others, Pulse{name, o, current.highlow})
            }
		case FLIPFLOP:
			others = processFlipflop(current, forward, states)
		case CONJUNCTION:
			others = processConjunction(current, forward, backward, conj_states)
		}
        for _, o := range others {
            if debug {
                fmt.Printf("%v -%v-> %v\n", name, current.highlow, o)
            }
            to_work = append(to_work, o)
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
	forward_links := make(map[string][]string)
	back_links := make(map[string][]string)
	modules_type := make(map[string]ModuleType)
	node_names := []string{}
	modules_state := make(map[string]bool)
	conjunction_states := make(map[string]map[string]bool)
	conjunction := make(map[string]bool)
	flipflops := make(map[string]bool)
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			fmt.Printf("%v \n", cleaned_line)
			links := strings.Split(cleaned_line, " -> ")
			source := replaceCharacters(links[0], "%&", "")
			destinations := links[1]
			forward_links[source] = []string{}
			back_links[source] = []string{}
			node_names = append(node_names, source)
			module_type := UNTYPED
			if IndexOfString(links[0], '%') >= 0 {
				module_type = FLIPFLOP
				flipflops[source] = true
			} else if IndexOfString(links[0], '&') >= 0 {
				module_type = CONJUNCTION
				conjunction[source] = true
			}
			modules_type[source] = module_type
			modules_state[source] = false
			conjunction_states[source] = make(map[string]bool)
			for _, d := range strings.Split(destinations, ", ") {
				// fmt.Printf("%v -> %v\n", source, d)
				forward_links[source] = append(forward_links[source], d)
			}
		}
	}
	node_names = removeDuplicates[string](node_names)
	for _, n := range node_names {
		for _, o := range node_names {
			if IndexOf[string](forward_links[o], n) >= 0 {
				back_links[n] = append(back_links[n], o)
			}
		}
		if modules_type[n] == CONJUNCTION {
			for _, s := range back_links[n] {
				conjunction_states[n][s] = false
			}
		}
	}

	if debug {
		fmt.Println("")
		fmt.Println("Run pulses")
	}

    found := false
	presses := 0
    to_watch := ""
    for _, name := range node_names {
        if IndexOf[string](forward_links[name], "rx")>=0 {
            to_watch = name
            break
        }
    }
    sources := back_links[to_watch]
    watch := make(map[Pulse]bool)
    for _, s := range sources {
        pulse := Pulse{s, to_watch, true}
        watch[pulse] = true
        fmt.Printf("Watching for pulse: %v\n", pulse)

    }
    seen := make(map[Pulse]bool)
    cycle := make(map[Pulse]int)
    len_sources := len(sources)
	for !found {
		presses++
        conj_observed := processSignal(presses, node_names, forward_links, back_links, modules_type, modules_state, conjunction_states)
        for _, pulse := range conj_observed {
            if watch[pulse] {
                fmt.Printf("%v: %v observed.\n", presses, pulse)
                if !seen[pulse] {
                    seen[pulse] = true
                    cycle[pulse] = presses
                    //fmt.Printf("%v\n", seen)
                    seen_count := 0
                    for _, b := range sources {
                        if seen[Pulse{b, to_watch, true}] {
                            seen_count++
                        }
                    }
                    // fmt.Printf("Seen count %v (len sources=%v)\n", seen_count, len_sources)
                    found = seen_count == len_sources
                    // fmt.Printf("Found = %v\n", found)
                }
            }
        }
	}

    result_cylces := []int{}
    for _, s := range sources {
        result_cylces = append(result_cylces, cycle[Pulse{s, to_watch, true}])
    }
    fmt.Printf("%v\n", result_cylces)
	// Collect results
	results := findLCM(result_cylces)

	fmt.Println("Final results:", results)

	fmt.Println("")
}
