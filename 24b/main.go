package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

func main() {
	if debug {
		fmt.Println("# Debug enabled")
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
	hailstones := []Hailstone{}
	fmt.Println("# Reading file", filename)
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			hailstones = append(hailstones, decodeHailstone(cleaned_line))
		}
	}

	if debug {
		fmt.Printf("# %v hailstones\n", len(hailstones))
	}
    equations := []Equation{}
    candidates := []map[string]Expr{}
    for i := 0; i < len(hailstones); i++ {
        t := Variable(fmt.Sprintf("t%d", i))
        h := hailstones[i]
        for _, comp := range []string{"x", "y", "z"} {
            var_pos := Variable(comp)
            var_vel := Multiplication{ []Expr{IntConstant(-1), Variable(fmt.Sprintf("v_%v", comp))} }
            pos_val := IntConstant(0)
            vel_val := IntConstant(0)
            switch comp {
            case "x":
                pos_val = IntConstant(-h.pos.x)
                vel_val = IntConstant(h.vel.x)
            case "y":
                pos_val = IntConstant(-h.pos.y)
                vel_val = IntConstant(h.vel.y)
            case "z":
                pos_val = IntConstant(-h.pos.z)
                vel_val = IntConstant(h.vel.z)
            }
            // hail_x + t * hail_vx = x + t * vx
            // <=> t * hail_vx - t * vx = x - hail_x
            // <=> t * (hail_vx - vx) = x - hail_x
            // <=> t = (x - hail_x) / (hail_vx - vx)
            upper := Addition{[]Expr{var_pos, pos_val}}
            lower := Addition{[]Expr{vel_val, var_vel}}
            eq := Equation{t, Fraction{upper, lower}}
            equations = append(equations, eq)
            fmt.Println(eq)

            new_candidate := make(map[string]Expr)
            new_candidate[comp] = Multiplication{[]Expr{IntConstant(-1), pos_val}}
            new_candidate["v_" + comp] = vel_val
            fmt.Println(new_candidate)
            candidates = append(candidates, new_candidate)
        }
    }
    fmt.Println(len(candidates))
    for i, cand := range candidates {
        fmt.Println(i, cand)
        for j, eq := range equations {
            if i != j {
                new_eq := eq.Eval(cand)
                fmt.Println(new_eq)
            }
        }
    }
}
