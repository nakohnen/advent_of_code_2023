package main

import (
	"flag"
	"fmt"
//	"math"
	"os"
	"strings"
	"sync"
)

const debug bool = true

type Position struct {
	x, y, z int
}

func addPositions(a, b Position) Position {
	return Position{x: a.x + b.x, y: a.y + b.y, z: a.z + b.z}
}

type Hailstone struct {
	pos Position
	vel Position
}

// var epsilon float64 = 0.001

func decodeHailstone(line string) Hailstone {
	at_split := strings.Split(line, "@")
	split_left := strings.Split(at_split[0], ",")
	split_right := strings.Split(at_split[1], ",")
	x := ToInt(strings.TrimSpace(split_left[0]))
	y := ToInt(strings.TrimSpace(split_left[1]))
	z := ToInt(strings.TrimSpace(split_left[2]))
	vx := ToInt(strings.TrimSpace(split_right[0]))
	vy := ToInt(strings.TrimSpace(split_right[1]))
	vz := ToInt(strings.TrimSpace(split_right[2]))
	return Hailstone{pos: Position{x: x, y: y, z: z}, vel: Position{x: vx, y: vy, z: vz}}
}

func detectCollision(v1, v2 Hailstone) (float64, float64, float64, float64, error) {
	dvx := v1.vel.x*v2.vel.y - v2.vel.x*v1.vel.y
	if dvx == 0 {
		return 0.0, 0.0, 0.0, 0.0, fmt.Errorf("No collision")
	}
	dx := v2.pos.x - v1.pos.x
	dy := v1.pos.y - v2.pos.y
	a := float64(dx*v2.vel.y+dy*v2.vel.x) / float64(dvx)
	b := float64(dx*v1.vel.y+dy*v1.vel.x) / float64(dvx)
	x1 := float64(v1.pos.x) + float64(v1.vel.x)*a
	y1 := float64(v1.pos.y) + float64(v1.vel.y)*a
	/*x2 := float64(v2.pos.x) + float64(v2.vel.x)*b
	y2 := float64(v2.pos.y) + float64(v2.vel.y)*b
	if math.Abs(x1-x2) > epsilon || math.Abs(y1-y2) > epsilon {
        fmt.Println("Collision is not at the same point, x1", x1, "y1", y1, "x2", x2, "y2", y2)
		fmt.Println("x1", x1, "y1", y1, "x2", x2, "y2", y2)
		panic("Collision is not at the same point")
	}*/
	return x1, y1, a, b, nil
}

func processWorker(to_work <-chan []Hailstone, wg *sync.WaitGroup, lower_bound, upper_bound int, results chan<- int) {
    result := 0
	for work := range to_work {
		if len(work) != 2 {
			panic("Invalid work")
		}
		v1 := work[0]
		v2 := work[1]
		x, y, a, b, err := detectCollision(v1, v2)
        if debug {
            fmt.Println("v1", v1, "v2", v2, "x", x, "y", y, "a", a, "b", b, "err", err)
        }
		if err != nil || a < 0 || b < 0 {
			continue
		}
		if x >= float64(lower_bound) && x <= float64(upper_bound) && y >= float64(lower_bound) && y <= float64(upper_bound) {
			result++
		}
	}
    results <- result
	wg.Done()
}

func main() {
	if debug {
		fmt.Println("Debug enabled")
	}
	cores_f := flag.Int("t", -1, "How many cores (threads) should we run?")
	filename_f := flag.String("f", "", "On which file should we run this?")
	lower_bound_f := flag.Int("l", 0, "Lower bound for test area")
	upper_bound_f := flag.Int("u", 0, "Upper bound for test area")
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

	lower_bound := *lower_bound_f
	upper_bound := *upper_bound_f
	if lower_bound > upper_bound || lower_bound == upper_bound {
		fmt.Println("Invalid bounds")
		os.Exit(1)
	}

	dat, err := os.ReadFile(filename)
	Check(err)

	// Read file
	hailstones := []Hailstone{}
	fmt.Println("Reading file", filename)
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			hailstones = append(hailstones, decodeHailstone(cleaned_line))
		}
	}

    if debug {
        fmt.Printf("%v hailstones\n", len(hailstones))
        fmt.Printf("Lower bound: %d, Upper bound: %d\n", lower_bound, upper_bound)
    }

	// Initialize workers
	numWorkers := cores
	var wg sync.WaitGroup
	if len(hailstones) < cores {
		numWorkers = len(hailstones)
	}
    fmt.Printf("Using %d workers\n", numWorkers)
    fmt.Printf("Using %d hailstones\n", len(hailstones))
	to_work := make(chan []Hailstone, len(hailstones))
	results := make(chan int, len(hailstones))
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processWorker(to_work, &wg, lower_bound, upper_bound, results)
	}

	// Send work to workers
	for i := 0; i < len(hailstones); i++ {
		for j := i + 1; j < len(hailstones); j++ {
			to_work <- []Hailstone{hailstones[i], hailstones[j]}
		}
	}
	close(to_work)

    // Wait for workers to finish
    go func() {
        wg.Wait()
        close(results)
    }()

	// Collect results
	collisions := 0
	for result := range results {
		collisions += result
	}
	fmt.Println("Collisions:", collisions)

	// Done
	fmt.Println("Done")
}
