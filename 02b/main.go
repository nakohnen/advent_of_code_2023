package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
func main() {
    dat, err := os.ReadFile("input.txt")
    check(err)

    text := string(dat)
    split_text := strings.Split(strings.TrimSpace(text), "\n")

    // Line by line
    sum := 0
    for _, line := range split_text {
        fmt.Println(line)
        line_split := strings.Split(line, ":")
        // We split each line into the left part game id
        // and the right part, the individual relevations
        // Left Part
        // game_id_split := strings.Split(line_split[0], " ")
        // game_id, err := strconv.Atoi(string(game_id_split[1]))
        // check(err)
        // Right Part
        games_split := strings.Split(line_split[1], ";")
        minimum_reds := 0
        minimum_greens := 0
        minimum_blues := 0
        for _, game := range games_split {
            for _, color := range strings.Split(game, ",") {
                color_split := strings.Split(strings.TrimSpace(color), " ")
                amount, err := strconv.Atoi(string(color_split[0]))
                check(err)
                switch color2 := color_split[1]; color2 {
                case "green": 
                    if amount > minimum_greens {
                        minimum_greens = amount
                    }
                case "red":
                    if amount > minimum_reds {
                        minimum_reds = amount
                    }
                case "blue":
                    if amount > minimum_blues {
                        minimum_blues = amount
                    }
                }
            }
        }
        power := minimum_greens * minimum_reds * minimum_blues
        fmt.Println(power)
        sum = sum + power 
    }
    fmt.Println(sum)
}
