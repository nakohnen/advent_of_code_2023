package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const debug bool = false

type Position struct {
	x, y int
}

func addPositions(a, b Position) Position {
	return Position{x: a.x + b.x, y: a.y + b.y}
}

type TileType int

const (
	WALL TileType = iota
	GROUND
	UP_SLOPE
	DOWN_SLOPE
	LEFT_SLOPE
	RIGHT_SLOPE
)

func decodeRuneToTile(r rune) TileType {
	switch r {
	case '#':
		return WALL
	case '.':
		return GROUND
	case '^':
		return UP_SLOPE
	case 'v':
		return DOWN_SLOPE
	case '<':
		return LEFT_SLOPE
	case '>':
		return RIGHT_SLOPE
	default:
		panic("Unknown tile type")
	}
}

func getNeighbours(tiles [][]TileType, pos Position) []Position {
	neighbours := []Position{}
	switch tiles[pos.y][pos.x] {
	case UP_SLOPE:
		neighbours = append(neighbours, Position{x: pos.x, y: pos.y - 1})
		return neighbours
	case DOWN_SLOPE:
		neighbours = append(neighbours, Position{x: pos.x, y: pos.y + 1})
		return neighbours
	case LEFT_SLOPE:
		neighbours = append(neighbours, Position{x: pos.x - 1, y: pos.y})
		return neighbours
	case RIGHT_SLOPE:
		neighbours = append(neighbours, Position{x: pos.x + 1, y: pos.y})
		return neighbours
	}

	if pos.x > 0 {
		ok := false
		switch tiles[pos.y][pos.x-1] {
		case GROUND:
			ok = true
		case LEFT_SLOPE:
			ok = true
		}
		if ok {
			neighbours = append(neighbours, Position{x: pos.x - 1, y: pos.y})
		}
	}
	if pos.x < len(tiles[0])-1 {
		ok := false
		switch tiles[pos.y][pos.x+1] {
		case GROUND:
			ok = true
		case RIGHT_SLOPE:
			ok = true
		}
		if ok {
			neighbours = append(neighbours, Position{x: pos.x + 1, y: pos.y})
		}
	}
	if pos.y > 0 {
		ok := false
		switch tiles[pos.y-1][pos.x] {
		case GROUND:
			ok = true
		case UP_SLOPE:
			ok = true
		}
		if ok {
			neighbours = append(neighbours, Position{x: pos.x, y: pos.y - 1})
		}
	}
	if pos.y < len(tiles)-1 {
		ok := false
		switch tiles[pos.y+1][pos.x] {
		case GROUND:
			ok = true
		case DOWN_SLOPE:
			ok = true
		}
		if ok {
			neighbours = append(neighbours, Position{x: pos.x, y: pos.y + 1})
		}
	}
	return neighbours
}

func isCrossroad(tiles [][]TileType, pos Position) bool {
    return tiles[pos.y][pos.x] != WALL && tiles[pos.y][pos.x] != GROUND
}

func isSlope(tiles [][]TileType, pos Position) bool {
    return tiles[pos.y][pos.x] == UP_SLOPE || tiles[pos.y][pos.x] == DOWN_SLOPE || tiles[pos.y][pos.x] == LEFT_SLOPE || tiles[pos.y][pos.x] == RIGHT_SLOPE
}

type PathComponent struct {
	id         int
	steps      int
	start, end Position
}

type Path struct {
	ids        []int
	components []PathComponent
	graph      map[int][]int
	start, end int
}

func walkPath(tiles [][]TileType, start Position, end Position) Path {
	path_starts := []Position{start}
	current_id := 1
	tiles_ids := make(map[Position]int)
	tiles_ids[start] = current_id
	paths := []PathComponent{}
	paths_graph := make(map[int][]int)
	paths_graph[current_id] = []int{}
	running_path_id := current_id
	path_ids := []int{current_id}
	start_id := current_id
	end_id := current_id
	// Find all paths
	for len(path_starts) > 0 {
		current_start := path_starts[0]
		current_id = tiles_ids[current_start]
		path_starts = path_starts[1:]
		if debug {
			fmt.Println("Walking from", current_start)
		}

		to_walk := []Position{current_start}
		walked := make(map[Position]bool)
		walked[current_start] = true
		steps := 0
		for len(to_walk) > 0 {
			pos := to_walk[0]
			to_walk = to_walk[1:]
			walked[pos] = true
			steps += 1

			if pos == end {
				// Found end
				end_id = current_id
			}

			// Current position is a slope
			is_slope := isSlope(tiles, pos)
			all_neighbours := getNeighbours(tiles, pos)
            neighbours := []Position{}
            for _, neighbour := range all_neighbours {
                if !walked[neighbour] {
                    neighbours = append(neighbours, neighbour)
                }
            }

			// Add current completed path to paths
			if is_slope || len(neighbours) != 1 || pos == end {
				if debug {
					fmt.Printf("Crossroad (or end) at %v, steps: %v for path id %v\n", pos, steps, current_id)
				}
				paths = append(paths, PathComponent{id: current_id, steps: steps, start: current_start, end: pos})
			}

			// Add neighbours to walk
            if len(neighbours) == 1 && !is_slope {
                to_walk = append(to_walk, neighbours[0])
            } else {
                if IndexOf[int](path_ids, current_id) == -1 {
				    paths = append(paths, PathComponent{id: current_id, steps: steps, start: current_start, end: pos})
                }
                for _, neighbour := range neighbours {
                    // Start new path if not already started
                    if tiles_ids[neighbour] == 0 {
                        // Create new path
                        running_path_id += 1
                        tiles_ids[neighbour] = running_path_id
                        paths_graph[running_path_id] = []int{}
                        path_starts = append(path_starts, neighbour)
                        path_ids = append(path_ids, running_path_id)
                        paths_graph[current_id] = append(paths_graph[current_id], running_path_id)
                        paths_graph[running_path_id] = append(paths_graph[running_path_id], current_id)
                        if debug {
                            fmt.Println("Starting new path", running_path_id, "with start from", neighbour)
                        }
                    } else {
                        // Connect to existing path
                        paths_graph[current_id] = append(paths_graph[current_id], tiles_ids[neighbour])
                        paths_graph[tiles_ids[neighbour]] = append(paths_graph[tiles_ids[neighbour]], current_id)
                        if debug {
                            fmt.Println("Connecting path", current_id, "to", tiles_ids[neighbour])
                        }
                    }
                }
            }
		}
	}
    for k, v := range paths_graph {
        paths_graph[k] = RemoveDuplicates[int](v)
    }

	return Path{ids: path_ids, components: paths, graph: paths_graph, start: start_id, end: end_id}
}

func getLongestPath(p Path) int {
	finished_paths := [][]int{}
	paths := [][]int{}
	paths = append(paths, []int{p.start})
	for len(paths) > 0 {
		current_path := paths[0]
		paths = paths[1:]
		current_id := current_path[len(current_path)-1]
		if current_id == p.end {
			finished_paths = append(finished_paths, current_path)
		} else {
			for _, next_id := range p.graph[current_id] {
                if IndexOf[int](current_path, next_id) == -1 {
                    new_path := []int{}
                    for _, id := range current_path {
                        new_path = append(new_path, id)
                    }
                    new_path = append(new_path, next_id)
                    paths = append(paths, new_path)
                }
			}
		}
	}
	paths = [][]int{}
	for _, path := range finished_paths {
		path_with_steps := []int{}
		for _, id := range path {
			found := false
			for _, comp := range p.components {
				if comp.id == id {
					path_with_steps = append(path_with_steps, comp.steps)
					found = true
					break
				}
			}
			if !found {
				panic("Could not find path component")
			}
		}
		paths = append(paths, path_with_steps)
	}

	if debug {
		fmt.Println("Finished paths:")
		for i, path := range finished_paths {
			fmt.Println(path, "->", paths[i])
		}
	}
	longest_path := 0
	for _, path := range paths {
        length := SumIntSlice(path) - 1
		longest_path = max(longest_path, length)
        if debug {
            fmt.Println("Path", path, "length", length)
        }
	}

	return longest_path
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
	tiles := [][]TileType{}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			row := []TileType{}
			for _, r := range cleaned_line {
				row = append(row, decodeRuneToTile(r))
			}
			tiles = append(tiles, row)
		}
	}
	width := len(tiles[0])
	height := len(tiles)
	start := Position{x: 1, y: 0}
	end := Position{x: width - 2, y: height - 1}
	// Walk path to find all paths
	path := walkPath(tiles, start, end)
	fmt.Println(path)

	// Find longest path
	longest_path := getLongestPath(path)
	fmt.Println("Longest path:", longest_path)
}
