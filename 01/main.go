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
    split_text := strings.Split(text, "\n")
    sum := 0
    for _, token := range split_text {
        number := 0
        for _, c := range token {
            if n, err := strconv.Atoi(string(c)); err == nil {
                number = n * 10
                break
            }
        }
        for _, c := range Reverse(token) {
            if n, err := strconv.Atoi(string(c)); err == nil {
                number = number + n
                break
            }
        }
        sum = sum + number
        fmt.Print(token)
        fmt.Print(" ")
        fmt.Print(number)
        fmt.Print(" ")
        fmt.Print(sum)
        fmt.Print("\n")
    }
}
