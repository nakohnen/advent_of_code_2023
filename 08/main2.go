package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
//	"math"
	//"sort"
	"strconv"
//	"sync"
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

func hasSuffix(line string, sub string) bool {
	if len(line) < len(sub) {
		return false
	}
	res := line[len(line)-len(sub):] == sub
	//fmt.Printf("line=%v suffix=%v sub=%v => %v\n", line, line[len(line)-len(sub):], sub, res)
	return res
}

func replaceCharacters(line string, rem string, rep string) string {
	result := line
	for _, c := range rem {
		result = strings.ReplaceAll(result, string(c), rep)
	}
	return result
}

func allEndIn(lines []string, suffix string) bool {
	for _, l := range lines {
		if hasSuffix(l, suffix) == false {
			return false
		}
	}
	return true
}

func oneEndIn(lines []string, suffix string) bool {
	for _, l := range lines {
		if hasSuffix(l, suffix) {
			return true
		}
	}
	return false
}

func concatInstructions(nodes []string, ins string) string {
	return strings.Join(nodes, "") + ins
}

func mapFunc(input string, nodeMap map[string][2]string, index int) string {
	return nodeMap[input][index]
}

func allEndNodes(nodes []string, end_nodes map[string]bool) bool {
	for _, n := range nodes {
		if !end_nodes[n] {
			return false
		}
	}
	return true
}

func findLoops(start_node string, instructions string, path map[string][2]string, end_nodes map[string]bool) (int, int) {
    fmt.Printf("Start node: %v with end_nodes %v\n", start_node, end_nodes)
	steps := 1
	counter := 0
	pre_loop := 0
	loop := 0

	cur_node := start_node
	end_node_visits := make(map[string]bool)
	for pre_loop == 0 || loop == 0 {
        next_index := 1
		if instructions[counter] == 'L' {
			next_index = 0
		}
		next_node := path[cur_node][next_index]

		sub_instructions := instructions[counter:]
//        fmt.Printf("-> %v with %v\n", next_node, sub_instructions)

		new_super_node := next_node + sub_instructions
		if end_nodes[next_node] {
            fmt.Printf("Visiting end_node %v with super node %v on steps%v\n", next_node, new_super_node, steps)
			if !end_node_visits[new_super_node] {
				end_node_visits[new_super_node] = true
			} else {
				if pre_loop == 0 {
					pre_loop = steps
				} else {
					loop = steps - pre_loop
					return pre_loop, loop
				}
			}
		}
        
        cur_node = next_node

		steps++
		counter++
		if counter >= len(instructions) {
			counter = 0
		}
	}
    return pre_loop, loop
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
	starting_nodes := []string{}
	end_nodes := make(map[string]bool)
	visits := make(map[string]int)
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
			visits[elements[0]] = 0
			if hasSuffix(elements[0], "A") {
				starting_nodes = append(starting_nodes, elements[0])
			}
			for _, s := range elements {
				if hasSuffix(s, "Z") {
					end_nodes[s] = true
				}
			}
			fmt.Printf("->%v\n", elements)
		}
	}

	cur_nodes := starting_nodes

	fmt.Printf("Starting nodes %v\n", cur_nodes)
	results := [][2]int{}

	for _, start_node := range cur_nodes {
		pre_loop, loop := findLoops(start_node, instructions, node_map, end_nodes)
		results = append(results, [2]int{pre_loop, loop})
	}

    steps := []int{}
    loops := []int{}

    for _, n := range results {
        steps = append(steps, n[0])
        loops = append(loops, n[1])
    }

    fmt.Printf("%v steps\n", steps) 
    fmt.Printf("%v loops\n", loops) 
    for !allElementsSame(steps) {
        lowest := steps[0]
        index := 0
        for i, s := range steps {
            if s < lowest {
                lowest = s
                index = i
            }
        }
        steps[index] += results[index][1]
    }



    fmt.Printf("%v steps\n", steps) 
	fmt.Printf("%v\n", results)

	fmt.Println("")
}
