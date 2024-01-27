package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

const debug bool = false

type Brick struct {
	low, high Position
}

func moveBrick(b Brick, vec Position) Brick {
	if debug {
		fmt.Printf("Moving brick %v by %v\n", b, vec)
	}
	return Brick{addPositions(b.low, vec), addPositions(b.high, vec)}
}

type Position struct {
	x, y, z int
}

func addPositions(a, b Position) Position {
	return Position{x: a.x + b.x, y: a.y + b.y, z: a.z + b.z}
}

func extractBrick(in string) Brick {
	tilde_split := strings.Split(in, "~")
	lower := strings.Split(tilde_split[0], ",")
	higher := strings.Split(tilde_split[1], ",")
	low := Position{ToInt(lower[0]), ToInt(lower[1]), ToInt(lower[2])}
	high := Position{ToInt(higher[0]), ToInt(higher[1]), ToInt(higher[2])}
	return Brick{low, high}
}

// Checks if a is on top of b
func onTopOf(a, b Brick) bool {
	if a.low.z == b.high.z+1 {
		for ax := a.low.x; ax <= a.high.x; ax++ {
			for ay := a.low.y; ay <= a.high.y; ay++ {
				if ax >= b.low.x && ax <= b.high.x && ay >= b.low.y && ay <= b.high.y {
					return true
				}
			}
		}
	}
	return false
}

func moveSingleBrick(bricks []Brick) bool {
	for i := 0; i < len(bricks); i++ {
		if bricks[i].low.z > 1 {
			supported := false
			for _, b := range bricks {
				if onTopOf(bricks[i], b) {
					if debug {
						fmt.Printf("Brick %v is supported by %v\n", bricks[i], b)
					}
					supported = true
					break
				}
			}
			if !supported {
				bricks[i] = moveBrick(bricks[i], Position{0, 0, -1})
				return true
			}
		}
	}
	return false
}

func getHighestBrickToFallOn(b Brick, supported_bricks []Brick) (int, Brick){
    highest_z := 0
    highest_b := Brick{}
    for b_x := b.low.x; b_x <= b.high.x; b_x++ {
        for b_y := b.low.y; b_y <= b.high.y; b_y++ {
            for _, b2 := range supported_bricks {
                if b2.low.x <= b_x && b_x <= b2.high.x&& b2.low.y <= b_y && b_y <= b2.high.y && b2.high.z > highest_z && b2.high.z < b.low.z {
                    highest_z = b2.high.z
                    highest_b = b2
                }
            }
        }
    }
    return highest_z, highest_b
}

// Move bricks down until they are supported
func letBricksFall(bricks []Brick) int {

    // Build list of bricks which can be used to check if a brick is supported
	supported_bricks := []Brick{}
    result := 0

    // For each brick, move it down until it is supported
	for i := 0; i < len(bricks); i++ {
		if bricks[i].low.z > 1 {
			supported := false
            // Move brick down until it is supported
            // Check for supported from back to front
            for j:=len(supported_bricks)-1; j>=0; j-- {
                b := supported_bricks[j]
                if onTopOf(bricks[i], b) {
                    if debug {
                        fmt.Printf("Brick %v is supported by %v\n", bricks[i], b)
                    }
                    supported = true
                    break
                }
            }
            // If not supported, move brick down
            if !supported {
                highest_z, highest_b := getHighestBrickToFallOn(bricks[i], supported_bricks)
                delta_z := bricks[i].low.z - highest_z - 1
                if debug {
                    fmt.Printf("Brick %v is not supported, moving down by %v (up to brick %v)", bricks[i], delta_z, highest_b)
                }
                bricks[i].low.z -= delta_z
                bricks[i].high.z -= delta_z
                result++
                if debug {
                    fmt.Printf(" -> %v\n", bricks[i])
                }
            }
            supported_bricks = append(supported_bricks, bricks[i])
		} else {
			supported_bricks = append(supported_bricks, bricks[i])
		}
	}

	// Assert  length of supported_bricks == length of bricks
	if len(supported_bricks) != len(bricks) {
		fmt.Printf("Error: letBricksFall failed %v %v\n", len(supported_bricks), len(bricks))
		os.Exit(1)
	}
    return result
}

// Checks if a brick can be safely removed
// supports: map of bricks to bricks which they support i.e. supports[b] are on top of b
// supported_by: map of bricks to bricks which support them i.e. b is on top of supported_by[b]
func canBeSafelyRemoved(b Brick, supports, supported_by map[Brick][]Brick) bool {
	for _, b2 := range supports[b] {
        // If b2 is only supported by b, then b cannot be safely removed
		if len(supported_by[b2]) == 1 {
			return false
		}
	}
	return true
}

func countFallingBricks(to_check_chan <-chan Brick, results_chan chan<- int, bricks []Brick, supports, supported_by map[Brick][]Brick, wg *sync.WaitGroup) {
    for b := range to_check_chan {
        to_work := []Brick{}
        for _, b2 := range bricks {
            if b2 != b {
                to_work = append(to_work, b2)
            }
        }
        results_chan <- letBricksFall(to_work)
    }
    wg.Done()
}

func main() {
	if debug {
		fmt.Println("Debug enabled")
	}
	cores_f := flag.Int("t", -1, "How many cores (threads) should we run?")
	filename_f := flag.String("f", "", "On which file should we run this?")
	flag.Parse()

	// The second element in os.Args is the first argument
	filename := *filename_f
	if len(filename) == 0 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}
	cores := *cores_f
	if cores < 0 {
		cores = 1
	}

	dat, err := os.ReadFile(filename)
	Check(err)

	// Read file
	fmt.Println("Reading file", filename)
	bricks := []Brick{}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			bricks = append(bricks, extractBrick(cleaned_line))
		}
	}
    fmt.Println("Read", len(bricks), "bricks")
	// Sort bricks
	fmt.Println("Sorting bricks")
    // Sort bricks by z, then y, then x
	sort.Slice(bricks, func(i, j int) bool {
		if bricks[i].low.z != bricks[j].low.z {
			return bricks[i].low.z < bricks[j].low.z
		} else if bricks[i].low.y != bricks[j].low.y {
			return bricks[i].low.y < bricks[j].low.y
		} else {
			return bricks[i].low.x < bricks[j].low.x
		}
	})

	fmt.Println("Moving bricks")
	letBricksFall(bricks)

	fmt.Println("Creating support graph")
    // supports: map of bricks to bricks which they support i.e. supports[b] are on top of b
    // supported_by: map of bricks to bricks which support them i.e. b is on top of supported_by[b]
	supported_by := make(map[Brick][]Brick)
	supports := make(map[Brick][]Brick)
	for _, b := range bricks {
		supported_by[b] = []Brick{}
		supports[b] = []Brick{}
	}
    // For each brick, check if it is supported by another brick
    // If so, add it to the support graph and the reverse graph
	for _, b := range bricks {
		for _, b2 := range bricks {
			if onTopOf(b, b2) {
				supports[b2] = append(supports[b2], b)
                supported_by[b] = append(supported_by[b], b2)
            } 
        }
    }

    to_check := []Brick{}
	for _, b := range bricks {
        if debug {
            fmt.Printf("Brick %v is supported by %v and supports %v ", b, supported_by[b], supports[b])
        }
		if canBeSafelyRemoved(b, supports, supported_by) {
			if debug {
				fmt.Printf("and can be safely removed\n")
			}
		} else {
            to_check = append(to_check, b)
            if debug {
                fmt.Printf("\n")
            }
        }
	}

    numWorkers := cores
    if numWorkers > len(to_check) {
        numWorkers = len(to_check)
    }
    fmt.Println("Using", numWorkers, "workers")
    fmt.Println("Checking", len(to_check), "bricks")
    fmt.Println("")

    // Create channels
    var wg sync.WaitGroup
    to_check_chan := make(chan Brick, len(to_check))
    results_chan := make(chan int, len(to_check))

    // Create workers
    for w:=0; w<numWorkers; w++ {
        wg.Add(1)
        go countFallingBricks(to_check_chan, results_chan, bricks, supports, supported_by, &wg)
    }

    // Send bricks to check
    for _, b := range to_check {
        to_check_chan <- b
    }
    close(to_check_chan)

    // Wait for workers to finish
    wg.Wait()

    // Collect results
    results := 0
    for i:=0; i<len(to_check); i++ {
        results += <-results_chan
    }
    close(results_chan)

	fmt.Println("Final results:", results)

	fmt.Println("")
}
