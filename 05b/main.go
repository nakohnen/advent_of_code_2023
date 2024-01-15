package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
	"math"
	"strconv"
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

func mapStep(mapping map[string][][]int, in_val []int, step string) []int {
	result := []int{}
	todo := []int{}
	for _, v := range in_val {
		todo = append(todo, v)
	}
	if val, ok := mapping[step]; ok {
		decode_slice := val
		for len(todo) > 0 {
			fmt.Printf("Step=%v, Todo=%v, Result=%v \n", step, todo, result)
			input_start := todo[0]
			input_end := todo[1]
			todo = todo[2:]
			if input_end < input_start {
				panic("Error with intervals")
			}
			intersection_found := false
			for _, interval := range decode_slice {
				source := interval[0]
				destination := interval[1]
				length := interval[2]

				source_end := source + length - 1

				destination_end := destination + length - 1

				if input_start == source && input_end == source_end {
					result = append(result, destination, destination_end)
					intersection_found = true
					break
				}

				if input_start <= source && source <= input_end {
					if input_start != source {
						todo = append(todo, input_start, source-1)
					}

					if source_end <= input_end {
						result = append(result, destination, destination_end)
						if source_end+1 <= input_end {
							todo = append(todo, source_end+1, input_end)
						}
					} else {
						// source_end > input_end, => we need to cut the the destination interval to the right length
						result = append(result, destination, destination+input_end-source)
					}
					intersection_found = true
					break
				}

				if source <= input_start && input_start <= source_end {
					if input_start == source {
						if input_end < source_end {
                            new_length := input_end - input_start

							result = append(result,
								destination,
								destination+new_length)
						} else {
							// input_end > source_end
							todo = append(todo, source_end+1, input_end)
							result = append(result, destination, destination_end)
						}
					} else if input_start == source_end {
						result = append(result, destination_end, destination_end)
						todo = append(todo, source_end+1, input_end)
					} else {
						// source < input_start && input_start < source_end
						if input_end < source_end {
							map_delta := destination - source
							result = append(result, input_start+map_delta, input_end+map_delta)
						} else if input_end > source_end {
							interval_delta := input_start - source
							todo = append(todo, source_end+1, input_end)
							result = append(result, destination+interval_delta, destination_end)
						} else {
							interval_delta := input_start - source
							result = append(result, destination+interval_delta, destination_end)
						}
					}

					intersection_found = true
					break
				}

			}
			if intersection_found == false {
				result = append(result, input_start, input_end)
			}
		}
	}
	return result
}

func main() {
	//dat, err := os.ReadFile("input_sample.txt")
	dat, err := os.ReadFile("input.txt")
	check(err)
	text := string(dat)

	min_location := math.MaxInt
	seeds := []int{}
	mode := ""
	var maps = map[string][][]int{}
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {

		if line == "" {
			//fmt.Println("")
			continue
		}
		fmt.Printf("%v => ", line)
		if hasPrefix(line, "seeds:") {
			seeds_str := strings.Split(line, " ")
			for _, seed := range seeds_str[1:] {
				val, err := strconv.Atoi(seed)
				check(err)
				seeds = append(seeds, val)
			}
			mode = "seeds"
			fmt.Printf("seeds=%v", seeds)
		} else if hasPrefix(line, "seed-to-soil") {
			mode = "soil"
			maps[mode] = [][]int{}
		} else if hasPrefix(line, "soil-to-fertilizer") {
			mode = "fertilizer"
			maps[mode] = [][]int{}
		} else if hasPrefix(line, "fertilizer-to-water") {
			mode = "water"
			maps[mode] = [][]int{}
		} else if hasPrefix(line, "water-to-light") {
			mode = "light"
			maps[mode] = [][]int{}
		} else if hasPrefix(line, "light-to-temperature") {
			mode = "temperature"
			maps[mode] = [][]int{}
		} else if hasPrefix(line, "temperature-to-humidity") {
			mode = "humidity"
			maps[mode] = [][]int{}
		} else if hasPrefix(line, "humidity-to-location") {
			mode = "location"
			maps[mode] = [][]int{}
		} else if mode != "" && line != "" {
			split := strings.Split(line, " ")
			destination_str := split[0]
			source_str := split[1]
			length_str := split[2]
			destination, err := strconv.Atoi(destination_str)
			check(err)
			source, err := strconv.Atoi(source_str)
			check(err)
			length, err := strconv.Atoi(length_str)
			check(err)

			maps[mode] = append(maps[mode], []int{source, destination, length})

			fmt.Printf("map=%v len=%v s=%v d=%v", mode, length, source, destination)
		} else {
			mode = ""

		}
		fmt.Println("")
	}
	flow := []string{"soil",
		"fertilizer",
		"water",
		"light",
		"temperature",
		"humidity",
		"location"}

	for i := 0; i < len(seeds); i = i + 2 {
		seed := seeds[i]
		amount := seeds[i+1]
		fmt.Printf("%v) seed range [%v, %v] length=%v \n", i/2, seed, seed+amount-1, amount)
		location := []int{seed, seed + amount - 1}
		for _, step := range flow {
			location = mapStep(maps, location, step)
			fmt.Printf("%v: location=%v\n", step, location)
		}

		for _, l := range location {
			if l < min_location {
				min_location = l
				fmt.Printf("=> min location=%v\n", min_location)

			}
		}
	}
	fmt.Printf("\n%v min location", min_location)
	fmt.Println("")
}
