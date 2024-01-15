package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	//"sync"

	//"unicode"
	//"math"
	//"math/big"

	//"sort"
	"strconv"
	//"sync"
)

const debug bool = false

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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func transposeStringSlice(slice []string) []string {
	rows := len(slice)
	cols := len(slice[0])
	result := []string{}
	if debug {
		fmt.Println("Transpose")
		fmt.Printf("Rows=%v cols=%v len(slice)=%v\n", rows, cols, len(slice))
	}
	var sb strings.Builder
	for col := 0; col < cols; col++ {
		sb.Reset()
		for row := 0; row < rows; row++ {
			sb.WriteByte(slice[row][col])
		}
		result = append(result, sb.String())
	}
	if debug {
		for i := 0; i < max(len(slice), len(result)); i++ {
			if i < len(slice) {
				fmt.Printf("%v ", slice[i])
			} else {
				sb.Reset()
				for j := 0; j < len(slice[0]); j++ {
					sb.WriteString(" ")
				}
				fmt.Printf("%v ", sb.String())
			}
			if i < len(result) {
				fmt.Printf("%v\n", result[i])
			} else {
				fmt.Println("")
			}
		}
	}

	return result
}

func countChar(s string, c rune) int {
	r := 0
	for _, sc := range s {
		if sc == c {
			r++
		}
	}
	return r
}

func sumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
}

func rotateRocks90Degrees(rocks []string) []string {
	if debug {
		fmt.Println("Rotate 90 degrees")
	}
	result := []string{}
	// col.row
	// 0.0 => 0.100
	// 100.0 => 0.0
	// 0.100 => 100.100
	// 100.100 => 100.0
	// ==> col1 becoms row1 but reversed
	cols := len(rocks[0])
	rows := len(rocks)
	var sb strings.Builder
	for i := 0; i < cols; i++ {
		sb.Reset()
		for j := rows - 1; j >= 0; j-- {
			sb.WriteByte(rocks[j][i])
		}
		result = append(result, sb.String())
	}
	if debug {
		for i, sub := range result {
			fmt.Printf("%v => %v\n", rocks[i], sub)
		}
	}
	return result
}

func rollRocks(rocks string) string {
	width := len(rocks)
	if debug {
		fmt.Printf("Roll rocks: %v with width %v\n", rocks, width)
	}
	sections := strings.Split(rocks, "#")
	var sb strings.Builder
	new_sections := []string{}
	for _, section := range sections {
		sb.Reset()
		len_section := len(section)
		stones := countChar(section, 'O')
		for i := 0; i < len_section; i++ {
			r := '.'
			if i < stones {
				r = 'O'
			}
			sb.WriteRune(r)
		}
		new_sections = append(new_sections, sb.String())
	}
	return strings.Join(new_sections, "#")
}

func calculateLoad(rocks string) int {
	width := len(rocks)
	if debug {
		fmt.Printf("Calculate load: %v with width=%v\n", rocks, width)
	}
	result := 0
	for i, r := range rocks {
		if r == 'O' {
			result += width - i
		}

	}
	return result
}

func runSubCycle(rocks []string) []string {
	// Roll north
	if debug {
		fmt.Println("Run Subcycle")
	}
	result := []string{}
	for _, sub := range transposeStringSlice(rocks) {
		res := rollRocks(sub)
		result = append(result, res)
	}
	result = transposeStringSlice(result)
	result = rotateRocks90Degrees(result)

	if debug {
		fmt.Println("Subcycle Result:")
		for i, sub := range result {
			fmt.Printf("%v => %v\n", rocks[i], sub)
		}
		fmt.Println("")
	}
	return result
}

func printRocks(rocks []string) {
	for _, r := range rocks {
		fmt.Println(r)
	}
	fmt.Println("")
}

func turnFullCylce(rocks []string) []string {
	result := runSubCycle(rocks)
	for i := 1; i < 4; i++ {
		result = runSubCycle(result)
	}
	return result
}

func main() {
	if debug {
		fmt.Println("Debug enabled")
	}
	filename_f := flag.String("f", "", "On which file should we run this?")
	flag.Parse()

	// The second element in os.Args is the first argument
	filename := *filename_f
	if len(filename) == 0 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

	dat, err := os.ReadFile(filename)
	check(err)

	// Read file
	image := []string{}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			image = append(image, cleaned_line)
			fmt.Printf("%v\n", cleaned_line)
		}
	}
	fmt.Println("")
	//fmt.Printf("Transpose image: width=%v height=%v\n", width, len(image))
	//image = transposeStringSlice(image)
	//**************************+
	repeats := make(map[string]int)
	repeats_set := make(map[string]bool)
	image_flat := strings.Join(image, "")

	repeats[image_flat] = 0
	repeats_set[image_flat] = true

	const max_turns int = 1000000000
	circle_found := false
	for turns := 1; turns <= max_turns; turns++ {
		image = turnFullCylce(image)
		image_flat = strings.Join(image, "")

		if repeats_set[image_flat] && !circle_found {
			circle := turns - repeats[image_flat]
			new_position := (max_turns - turns) % circle
			turns = max_turns - new_position
			circle_found = true
		} else {
			repeats[image_flat] = turns
			repeats_set[image_flat] = true
		}

		fmt.Printf("Cycle: %v\n", turns)
	}

	// Collect results
	results := 0
	for _, row := range transposeStringSlice(image) {
		results += calculateLoad(row)
	}

	printRocks(image)

	fmt.Println("Final results:", results)

	fmt.Println("")
}
