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
	// dat, err := os.ReadFile("input_sample.txt")
	dat, err := os.ReadFile("input.txt")
	check(err)
	text := string(dat)

	sum := 0
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
		score := 0
		for _, n := range played {
            if winning_set[n] {
				if score == 0 {
					score = 1
				} else {
					score = score * 2
				}
			}
		}
		fmt.Printf("%v  score=%v \n", len(winning), score)
		sum += score
	}
	fmt.Println(sum)
}
