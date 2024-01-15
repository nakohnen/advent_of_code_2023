package main

import (
	"fmt"
	"os"
	"strings"
	//"unicode"
	//	"math"
	"sort"
	"strconv"
	//"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var cardStrength string = "AKQJT98765432"
var typeStrength = map[string]int{
	"1":  1,
	"2":  2,
	"22": 3,
	"3":  4,
	"fh": 5,
	"4":  6,
	"5":  7,
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

func getCardCount(hand string, card rune) int {
	result := 0
	for _, h := range hand {
		if h == card {
			result++
		}
	}
	return result
}

func getHandType(strengthCount []int) string {
	triplets := 0
	pairs := 0
	ones := 0

	for _, c := range strengthCount {
		switch c {
		case 5:
			return "5"
		case 4:
			return "4"
		case 3:
			triplets++
		case 2:
			pairs++
		case 1:
			ones++
		}
	}

	if triplets > 0 {
		if pairs > 0 {
			return "fh"
		}
		return "3"
	}
	if pairs > 0 {
		if pairs > 1 {
			return "22"
		}
		return "2"
	}
	return "1"
}

func compareHands(hand string, other string, strength string, typeStrength map[string]int) int {
	if len(strength) == 0 {
		return -1
	}
	count_type_1 := []int{}
	count_type_2 := []int{}
	for i := 0; i < len(strength); i++ {
		count_type_1 = append(count_type_1, 0)
		count_type_2 = append(count_type_2, 0)
	}
    //fmt.Println("Here")

	for i, c := range strength {
		count_type_1[i] = getCardCount(hand, c)
		count_type_2[i] = getCardCount(other, c)
	}
    //fmt.Printf("Hand 1: %v\n", count_type_1)
    //fmt.Printf("Hand 2: %v\n", count_type_2)

	hand1 := typeStrength[getHandType(count_type_1)]
	hand2 := typeStrength[getHandType(count_type_2)]
    //fmt.Printf("%v vs %v <=> %v vs %v\n", hand1, hand2, 
    //getHandType(count_type_1),
    //getHandType(count_type_2))
	if hand1 > hand2 {
		return 1
	} else if hand1 < hand2 {
		return 0
	} else {
		for i := 0; i < len(hand); i++ {
			card1 := strings.Index(strength, string(hand[i]))
			card2 := strings.Index(strength, string(other[i]))
			if card1 < card2 {
				return 1
			} else if card1 > card2 {
				return 0
			}
		}
		return -1
	}

}

type Hands []string

func (h Hands) Len() int {
	return len(h)
}

func (h Hands) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h Hands) Less(i, j int) bool {
    res :=  compareHands(h[i], h[j], cardStrength, typeStrength) == 0
    fmt.Printf("%v<%v==%v\n", h[i], h[j], res) 

    return res
}

func main() {
	//dat, err := os.ReadFile("input_sample.txt")
	dat, err := os.ReadFile("input.txt")
	check(err)
	hands := []string{}
	bids := make(map[string]int)
	for _, line := range strings.Split(strings.TrimSpace(string(dat)), "\n") {
		fmt.Printf("%v => ", line)
		line_split := strings.Split(line, " ")
		hand := line_split[0]
		bid := toInt(line_split[1])
		bids[hand] = bid

		hands = append(hands, hand)
		fmt.Printf("hand=%v bid=%v\n", hand, bid)
	}
	fmt.Printf("%v\n", hands)
	sort.Sort(Hands(hands))
	fmt.Printf("%v\n", hands)
    sum := 0
    for r, h := range(hands) {
        sum += (r+1) * bids[h]
        fmt.Printf("rank=%v, hand=%v, sum=%v\n", r+1, h, sum)
    }

	fmt.Println("")
}
