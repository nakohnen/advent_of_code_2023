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

func getSubSequence(seq []int) []int {
	if len(seq) <= 1 {
		return []int{}
	}
	new_seq := []int{}
	for i := 0; i < len(seq)-1; i++ {
		new_seq = append(new_seq, seq[i+1]-seq[i])
	}
	return new_seq
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
	sum := 0
	for _, line := range strings.Split(strings.TrimSpace(string(dat)), "\n") {
		seq := []int{}
		for _, i := range strings.Split(line, " ") {
			seq = append(seq, toInt(i))
		}
		fmt.Printf("%v ->\n", seq)

		predictor_seq := []int{}

		for !allElementsSame(seq) {
			predictor_seq = append(predictor_seq, seq[0])
			seq = getSubSequence(seq)
			fmt.Printf("Subseq=%v\n", seq)
		}
		if len(seq) > 0 {
			predictor_seq = append(predictor_seq, seq[0])
		}
        fmt.Printf("predictor_seq=%v\n", predictor_seq)
		sub_sum := 0
        for i:=len(predictor_seq)-1;i>=0;i-- {
			sub_sum = predictor_seq[i] - sub_sum
            fmt.Printf("subsum=%v at %v\n", sub_sum, i)
		}

		sum += sub_sum
		fmt.Printf("(%v)  => sum=%v\n", sub_sum, sum)

	}
	fmt.Printf("sum=%v\n", sum)
	fmt.Println("")
}
