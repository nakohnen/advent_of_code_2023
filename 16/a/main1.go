package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	//	"sync"

	"unicode"
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

func sumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
}

func hashString(s string) int {
	result := 0
	for _, r := range s {
		if unicode.IsPrint(r) {
			//fmt.Printf("%v = %v\n", r, string(r))
			result += int(r)
			result *= 17
			result %= 256
		}
	}

	return result
}

func IndexOf[T comparable](s []T, element T) int {
	for i, e := range s {
		if e == element {
			return i
		}
	}
	return -1
}

func IndexOfString(s string, r rune) int {
	for i, r2 := range s {
		if r2 == r {
			return i
		}
	}
	return -1
}

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

type Position struct {
	x, y int
}

type Mirror int

const (
    GROUND Mirror = iota
	HORIZONTAL
	VERTICAL
	DIAGONAL_NESW
	DIAGONAL_NWSE
)

type Beam struct {
	pos Position
	dir Direction // Incoming direction which should be unique
}

type BeamSet struct {
	set       map[Beam]bool
	container []Beam
}

func addToSet(bs *BeamSet, b Beam) bool {
	if bs.set[b] {
		return false
	}
	bs.set[b] = true
	bs.container = append(bs.container, b)
	return true
}

func resolveLazer(grid map[Position]Mirror, inc_b Beam) []Beam {
	mirror := grid[inc_b.pos]
	out_beams := []Beam{}
	switch dir := inc_b.dir; dir {
	case UP:
		switch mirror {
		case HORIZONTAL:
			split_left := Beam{Position{x: inc_b.pos.x - 1, y: inc_b.pos.y}, LEFT}
			out_beams = append(out_beams, split_left)
			split_right := Beam{Position{x: inc_b.pos.x + 1, y: inc_b.pos.y}, RIGHT}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NESW:
			split_right := Beam{Position{x: inc_b.pos.x + 1, y: inc_b.pos.y}, RIGHT}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NWSE:
			split_left := Beam{Position{x: inc_b.pos.x - 1, y: inc_b.pos.y}, LEFT}
			out_beams = append(out_beams, split_left)
		default:
			out_beam := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y - 1}, UP}
			out_beams = append(out_beams, out_beam)
		}
	case DOWN:
		switch mirror {
		case HORIZONTAL:
			split_left := Beam{Position{x: inc_b.pos.x - 1, y: inc_b.pos.y}, LEFT}
			out_beams = append(out_beams, split_left)
			split_right := Beam{Position{x: inc_b.pos.x + 1, y: inc_b.pos.y}, RIGHT}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NESW:
			split_right := Beam{Position{x: inc_b.pos.x - 1, y: inc_b.pos.y}, LEFT}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NWSE:
			split_left := Beam{Position{x: inc_b.pos.x + 1, y: inc_b.pos.y}, RIGHT}
			out_beams = append(out_beams, split_left)
		default:
			out_beam := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y + 1}, DOWN}
			out_beams = append(out_beams, out_beam)
		}
	case LEFT:
		switch mirror {
		case VERTICAL:
			split_left := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y - 1}, UP}
			out_beams = append(out_beams, split_left)
			split_right := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y + 1}, DOWN}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NESW:
			split_right := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y + 1}, DOWN}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NWSE:
			split_left := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y - 1}, UP}
			out_beams = append(out_beams, split_left)
		default:
			out_beam := Beam{Position{x: inc_b.pos.x - 1, y: inc_b.pos.y}, LEFT}
			out_beams = append(out_beams, out_beam)
		}
	case RIGHT:
		switch mirror {
		case VERTICAL:
			split_left := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y - 1}, UP}
			out_beams = append(out_beams, split_left)
			split_right := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y + 1}, DOWN}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NWSE:
			split_right := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y + 1}, DOWN}
			out_beams = append(out_beams, split_right)
		case DIAGONAL_NESW:
			split_left := Beam{Position{x: inc_b.pos.x, y: inc_b.pos.y - 1}, UP}
			out_beams = append(out_beams, split_left)
		default:
			out_beam := Beam{Position{x: inc_b.pos.x + 1, y: inc_b.pos.y}, RIGHT}
			out_beams = append(out_beams, out_beam)
		}
	}
	if debug {
		fmt.Printf("%v %v => %v\n", inc_b, mirror, out_beams)
	}
	return out_beams
}

func paintGrid(grid map[Position]Mirror, height, width int, bs BeamSet, mirror_decode map[Mirror]rune) {
    var sb strings.Builder
    for y:=0;y<height;y++ {
        sb.Reset()
        for x:=0;x<width;x++ {
            pos := Position{x, y}
            tile :=  mirror_decode[grid[pos]]
            found := false
            for _, beam := range bs.container {
                if beam.pos == pos {
                    found = true
                    break
                }
            }
            if found {
                tile = '#'
            }
            sb.WriteRune(tile)
        }
        fmt.Println(sb.String())
    }
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
	decode_grid := make(map[rune]Mirror)
    reverse_decode := make(map[Mirror]rune)
	decode_grid['|'] = VERTICAL
    reverse_decode[VERTICAL] = '|'
	decode_grid['-'] = HORIZONTAL
    reverse_decode[HORIZONTAL] = '-'
	decode_grid['.'] = GROUND
    reverse_decode[GROUND] = '.'
	decode_grid['/'] = DIAGONAL_NESW
    reverse_decode[DIAGONAL_NESW] = '/'
	decode_grid['\\'] = DIAGONAL_NWSE
    reverse_decode[DIAGONAL_NWSE] = '\\'
	beam_set := BeamSet{
		set:       make(map[Beam]bool),
		container: []Beam{},
	}
	grid := make(map[Position]Mirror)
	height := 0
	width := 0
	for y, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
        if len(cleaned_line) > 0 {
		for x, r := range cleaned_line {
			position := Position{x, y}
			grid[position] = decode_grid[r]
			if x > width {
				width = x
			}
		}
		if y > height {
			height = y
		}
		if debug {
			fmt.Printf("%v %v \n", cleaned_line, len(cleaned_line))
		}
    }
	}
	height++
	width++
	//**************************+

	working_slice := []Beam{{Position{0, 0}, RIGHT}}
    addToSet(&beam_set, working_slice[0])
	for len(working_slice) > 0 {
		beam := working_slice[0]
		working_slice = working_slice[1:]

		for _, out_beam := range resolveLazer(grid, beam) {
			pos := out_beam.pos
			if pos.x >= 0 && pos.y >= 0 && pos.x < width && pos.y < height {
				if addToSet(&beam_set, out_beam) {
					working_slice = append(working_slice, out_beam)
				}
			}
		}
        if debug {
            paintGrid(grid, height, width, beam_set, reverse_decode)
        }
	}

	pos_set := make(map[Position]bool)

	results := 0
	for _, beam := range beam_set.container {
		if !pos_set[beam.pos] {
			results++
			pos_set[beam.pos] = true
		}
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
