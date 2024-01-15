package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
	//	"math"
	//"sort"
	"strconv"
	//"sync"
)

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

func main() {
	// Check if an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

	// The second element in os.Args is the first argument
	filename := os.Args[1]

	//dat, err := os.ReadFile("input_sample.txt")
	dat, err := os.ReadFile(filename)
	//dat, err := os.ReadFile("input.txt")
	check(err)

	node_map := make(map[string][2]string)
	instructions := ""
	for n, line := range strings.Split(strings.TrimSpace(string(dat)), "\n") {
		if n == 0 {
			instructions = line
			fmt.Printf("instructions %v\n", len(instructions))
		} else if n >= 2 {
			// ZZZ = (ZZZ, ZZZ)
			edited := line
			edited = replaceCharacters(edited, " ()", "")
			edited = replaceCharacters(edited, ",", "=")
			elements := strings.Split(edited, "=")
			node_map[elements[0]] = [2]string{elements[1], elements[2]}
			fmt.Printf("->%v\n", elements)
		}
	}

	cur_node := "AAA"
	target := "ZZZ"
	instruction_counter := 0
	steps := 0
	for cur_node != target {
		next_index := 1
		if instruction_counter >= len(instructions) {
			instruction_counter = 0
		}
		if instructions[instruction_counter] == 'L' {
			next_index = 0
		}
		cur_node = node_map[cur_node][next_index]
		instruction_counter++
		steps++
	}
	fmt.Printf("steps=%v\n", steps)
	fmt.Println("")
}
