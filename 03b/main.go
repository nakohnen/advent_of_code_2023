package main

import (
    "fmt"
    "os"
    "strings"
    "unicode"
    "math"
    "strconv"
)

type grid_value struct {
    value    int
    x_min    int
    x_max    int
    y        int
}

type grid_symbol struct {
    x      int
    y      int
    values []grid_value
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func AddSymbol(x int, y int, symbol_list []grid_symbol ) []grid_symbol {
    s := grid_symbol{x: x, y: y}
    return append(symbol_list, s)
}

func AddValue (value_text string, x_min int, x_max int, y int, value_list []grid_value) []grid_value {
    value, err := strconv.Atoi(value_text)
    check(err)
    v := grid_value{value: value, 
            x_min: x_min, 
            x_max: x_max, 
            y: y}
    return append(value_list, v)
}

func main() {
    //dat, err := os.ReadFile("input_sample.txt")
    dat, err := os.ReadFile("input.txt")
    check(err)

    text := string(dat)
    split_text := strings.Split(strings.TrimSpace(text), "\n")

    // Line by line
    symbols := []grid_symbol{}
    values := []grid_value{}


    sum := 0
    for y, line := range split_text {
        fmt.Printf("%v: %v ", y, line)
        var buffer strings.Builder
        x_min := math.MaxInt
        x_max := 0
        for x, token := range line {
            if token == '.' {
                // We encountered empty text
                if buffer.String() != "" {
                    x_max = x - 1
                    values = AddValue(buffer.String(), x_min, x_max, y, values)
                    fmt.Printf("value %v (%v, %v),", buffer.String(), x_min, y)
                }
                buffer.Reset()
                x_min = math.MaxInt
            } else if unicode.IsDigit(token) {
                // We encountered a digit
                _, err = buffer.WriteString(string(token))
                check(err)
                x_min = func(a, b int) int { if a < b { return a } else { return b } }(x_min, x)
            } else {
                // We encountered a symbol
                if buffer.String() != "" {
                    x_max = x - 1
                    values = AddValue(buffer.String(), x_min, x_max, y, values)
                    fmt.Printf("value %v (%v, %v),", buffer.String(), x_min, y)
                }
                if token == '*' {
                    symbols = AddSymbol(x, y, symbols)
                    fmt.Printf("symbol %v (%v, %v),", string(token), x, y)
                }
                
                buffer.Reset()
                x_min = math.MaxInt
            }
            x_max = x
                
        }
        if buffer.String() != "" {
            values = AddValue(buffer.String(), x_min, x_max, y, values)
            fmt.Printf("value %v (%v, %v),", buffer.String(), x_min, y)
        }
        fmt.Println("")
    }
    
    for i, s := range symbols {
        for _, v := range values {
            found := false
            for x := v.x_min; x <= v.x_max; x++ {
                y := v.y
                if s.x -1 <= x && x <= s.x + 1 && s.y-1 <= y && y <= s.y+1 {
                    found = true
                    break
                }
            }
            if found {
                fmt.Printf("value %v is near symbol %v\n", v, s)
                symbols[i].values = append(symbols[i].values, v) 
            }
        }
    }

    for _, s := range symbols {
        if len(s.values) == 2 {
            sum = sum + s.values[0].value * s.values[1].value
        }
    }
    fmt.Println(sum)
}
