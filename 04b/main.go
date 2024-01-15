package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
	//"math"
	// "strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	////dat, err := os.ReadFile("input_sample.txt")
	dat, err := os.ReadFile("input.txt")
	check(err)
	text := string(dat)

	sum := 0
	won := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		fmt.Printf("%v => ", line)
		line_split := strings.Split(line, ":")
		card_split := strings.Split(line_split[1], "|")
		winning := strings.Split(card_split[0], " ")
		played := strings.Split(card_split[1], " ")
		var winning_set map[string]bool = make(map[string]bool)
		for _, n := range winning {
			if n != "" {
				winning_set[n] = true
				fmt.Printf("%v ", strings.TrimSpace(string(n)))
			}
		}
		multiplier := 1 + won[0]
		for i := 0; i < 9; i++ {
			won[i] = won[i+1]
		}
		won[9] = 0

		score := 0
		for _, n := range played {
			if winning_set[n] {
				score = score + 1
			}
		}
		for i := 0; i < score; i++ {
			won[i] = won[i] + multiplier
		}
		sum += multiplier
		fmt.Printf("%v  score=%v sum=%v \n", len(winning), score, sum)
	}
	fmt.Println(sum)
}
