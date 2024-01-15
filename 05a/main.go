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

func mapStep(mapping map[string][][]int, in_val int, step string) int {
	if val, ok := mapping[step]; ok {
		decode_slice := val
		for _, interval := range decode_slice {
			if in_val >= interval[0] && in_val <= interval[0]+interval[2] {
				delta := in_val - interval[0]
				return interval[1] + delta
			}
		}
	}
	return in_val
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

	for _, seed := range seeds {
		location := seed
		for _, step := range flow {
			location = mapStep(maps, location, step)
			//if val, ok := maps[step][location]; ok {
			//	location = val
			//}
		}
		if location < min_location {
			min_location = location
		}
	}
	fmt.Printf("%v min location\n", min_location)
	fmt.Println("End")
}
