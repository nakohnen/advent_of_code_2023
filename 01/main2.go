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
func ReplaceNumbers(s string) string {
    s2 := strings.Clone(s)
    s2 = strings.Replace(s2, "one", "1", -1)
    s2 = strings.Replace(s2, "two", "2", -1)
    s2 = strings.Replace(s2, "three", "3", -1)
    s2 = strings.Replace(s2, "four", "4", -1)
    s2 = strings.Replace(s2, "five", "5", -1)
    s2 = strings.Replace(s2, "six", "6", -1)
    s2 = strings.Replace(s2, "seven", "7", -1)
    s2 = strings.Replace(s2, "eight", "8", -1)
    s2 = strings.Replace(s2, "nine", "9", -1)
    return s2

}
func main() {
    dat, err := os.ReadFile("input2.txt")
    check(err)

    text := string(dat)
    split_text := strings.Split(text, "\n")
    sum := 0
    for _, token := range split_text {
        number := 0

        for i := 0; i < len(token); i++ {
            sub_s := token[0:i+1]
            sub_s = ReplaceNumbers(sub_s)
            fmt.Print(">")
            fmt.Println(sub_s)
            if n, err := strconv.Atoi(string(sub_s[len(sub_s)-1])); err == nil {
                number = n * 10
                break
            }
        }

        for i := 0; i < len(token); i++ {
            position := len(token)-i
            sub_s := token[position-1:len(token)]
            sub_s = ReplaceNumbers(sub_s)
            fmt.Print("<")
            fmt.Println(sub_s)
            if n, err := strconv.Atoi(string(sub_s[0])); err == nil {
                number = number + n
                break
            }
        }

        sum = sum + number
        fmt.Print("> ")
        fmt.Print(token)
        fmt.Print(" ")
        fmt.Print(number)
        fmt.Print(" ")
        fmt.Print(sum)
        fmt.Print("\n")
    }
}
