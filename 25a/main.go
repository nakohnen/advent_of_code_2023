package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
)

const debug bool = false

func countNodeSize(n string, delimeter string) int {
    return len(strings.Split(n, delimeter))
}


func mergeNodes(graph map[string][]string, weights map[string]map[string]int, nodes []string, merge1, merge2 string) (map[string][]string, map[string]map[string]int, []string){
    return mergeNodesGroup(graph, weights, nodes, []string{merge1, merge2})
}

func mergeNodesGroup(graph map[string][]string, weights map[string]map[string]int, nodes []string, merge []string) (map[string][]string, map[string]map[string]int, []string){
    if len(merge) == 0 {
        return graph, weights, nodes
    }
    var sb strings.Builder
    sb.WriteString(merge[0])
    for _, node := range merge[1:] {
        sb.WriteString("-")
        sb.WriteString(node)
    }
    to_merge := sb.String()


    new_nodes := []string{}
    for _, node := range nodes {
        // If the node is not in the merge list
        if IndexOf[string](merge, node) == -1 {
            if debug {
                fmt.Println("Is not in merge list", node, "-> add to new nodes")
            }
            new_nodes = append(new_nodes, node)
        }
    }
    // Add the new node
    new_nodes = append(new_nodes, to_merge)

    if debug {
        fmt.Println("Merging", merge)
        fmt.Println("New nodes", new_nodes)
        fmt.Println("Old nodes", nodes)
    }

    // Create new graph
    new_graph := make(map[string][]string)
    new_weights := make(map[string]map[string]int)
    for _, node := range new_nodes {
        new_graph[node] = []string{}
        new_weights[node] = make(map[string]int)
        for _, other_node := range new_nodes {
            new_weights[node][other_node] = 0
        }
    }

    for _, node := range nodes {
        // If the node is not in the merge list
        if IndexOf[string](merge, node) == -1 {
            // For each neighbour of the node copy the old weights or 
            // recalculate the weights to the merged node and do it vice versa
            // so the merged node also has the correct weights and points to 
            // the old node.
            for _, neighbour := range graph[node] {

                // If the neighbour is in the merge list
                if IndexOf[string](merge, neighbour) != -1 {
                    if debug {
                        fmt.Println("Pointing to merged node", node, neighbour, "=>", to_merge)
                    }

                    // Point to the new node and back
                    new_graph[node] = append(new_graph[node], to_merge)
                    new_graph[to_merge] = append(new_graph[to_merge], node)

                    // Add weights based on the combined weights
                    new_weight := 0
                    for _, merge_node := range merge {
                        new_weight += weights[node][merge_node]
                    }
                    new_weights[node][to_merge] = new_weight
                    new_weights[to_merge][node] = new_weight 
                } else {
                    // Neighbour is not in the merge list
                    // Copy the old weights
                    new_graph[node] = append(new_graph[node], neighbour)
                    if debug {
                        fmt.Println("Copying weight", node, neighbour)
                        fmt.Println("Old weight", weights[node][neighbour])
                        fmt.Println("New nodes", new_nodes)
                        fmt.Println("Old weights", weights)
                    }
                    new_weights[node][neighbour] = weights[node][neighbour]
                }
            }
        }
    }
    // Do some housekeeping and remove duplicates
    // As we dont keep track of duplicates when merging in the code above
    for _, node := range new_nodes {
        new_graph[node] = RemoveDuplicates[string](new_graph[node])
    }

    return new_graph, new_weights, new_nodes
}

func listAllPossibleNeighbours(graph map[string][]string, nodes []string, ignore []string) []string {
    // Create a list of all possible neighbours which are not in nodes and ignore
    neighbours := []string{}
    for _, node := range nodes {
        for _, neighbour := range graph[node] {
            if IndexOf[string](nodes, neighbour) == -1 && IndexOf[string](ignore, neighbour) == -1 && IndexOf[string](neighbours, neighbour) == -1 {
                neighbours = append(neighbours, neighbour)
            }
        }
    }

    for _, neighbour := range neighbours {
        if IndexOf[string](nodes, neighbour) != -1 {
            panic("Neighbour is in nodes")
        }
        if IndexOf[string](ignore, neighbour) != -1 {
            panic("Neighbour is in ignore")
        }
    }

    return neighbours
}

func mergeNodesUntil2Remain(graph map[string][]string, weights map[string]map[string]int, nodes []string) (int, int){
    merged_1 := []string{}
    merged_2 := []string{}

    first := rand.Intn(len(nodes))
    second := rand.Intn(len(nodes))
    for second == first {
        second = rand.Intn(len(nodes))
    }
    merged_1 = append(merged_1, nodes[first])
    merged_2 = append(merged_2, nodes[second])

    run_nodes := []string{}
    for _, node := range nodes {
        // If the node is in neither of the merged lists
        if IndexOf(merged_1, node) == -1 && IndexOf(merged_2, node) == -1 {
            run_nodes = append(run_nodes, node)
        }
    }

    for len(run_nodes) > 0 {
        first_neighbours := listAllPossibleNeighbours(graph, merged_1, merged_2)
        if len(first_neighbours) > 0 {
            first := rand.Intn(len(first_neighbours))
            to_merge_1 := first_neighbours[first]
            merged_1 = append(merged_1, to_merge_1)
        }

        second_neighbours := listAllPossibleNeighbours(graph, merged_2, merged_1)
        if len(second_neighbours) > 0 {
            second := rand.Intn(len(second_neighbours))
            to_merge_2 := second_neighbours[second]
            merged_2 = append(merged_2, to_merge_2)
        }
        
        run_nodes = []string{}
        for _, node := range nodes {
            // If the node is in neither of the merged lists
            if IndexOf(merged_1, node) == -1 && IndexOf(merged_2, node) == -1 {
                run_nodes = append(run_nodes, node)
            }
        }
    }

    size_1 := len(merged_1)
    size_2 := len(merged_2)
    cut := 0
    for _, node_1 := range merged_1 {
        for _, node_2 := range merged_2 {
            cut += weights[node_1][node_2]
        }
    }
    if size_1 + size_2 != len(nodes) {
        // Panic, the sizes should add up to the number of nodes
        // Print the nodes and the merged lists
        fmt.Println("Nodes", nodes)
        fmt.Println("Merged 1", merged_1)
        fmt.Println("Merged 2", merged_2)
        panic("Size of merged lists does not add up to the number of nodes")
    }

    return cut, size_1 * size_2
}

func KargersRun(graph map[string][]string, weights map[string]map[string]int, nodes []string) (int, int) {

        run_nodes := []string{}
        for _, node := range nodes {
            run_nodes = append(run_nodes, node)
        }
        run_graph := make(map[string][]string)
        run_weights := make(map[string]map[string]int)
        for _, node := range run_nodes {
            run_graph[node] = []string{}
            run_weights[node] = make(map[string]int)
            for _, other_node := range run_nodes {
                run_weights[node][other_node] = 0
            }
            for _, neighbour := range graph[node] {
                run_graph[node] = append(run_graph[node], neighbour)
                run_weights[node][neighbour] = weights[node][neighbour]
            }
        }

        for len(run_nodes) > 2 {
            first := rand.Intn(len(run_nodes))
            second := rand.Intn(len(run_nodes))
            for second == first {
                second = rand.Intn(len(run_nodes))
            }
            merge1 := run_nodes[first]
            merge2 := run_nodes[second]
            run_graph, run_weights, run_nodes = mergeNodes(run_graph, run_weights, run_nodes, merge1, merge2)

        }

        group1_size := countNodeSize(run_nodes[0], "-")
        group2_size := countNodeSize(run_nodes[1], "-")

        cuts := run_weights[run_nodes[0]][run_nodes[1]]
        groups_size_multiplied := group1_size * group2_size
        return cuts, groups_size_multiplied
}

func KargersRunAlt(graph map[string][]string, weights map[string]map[string]int, nodes []string) (int, int) {
    return mergeNodesUntil2Remain(graph, weights, nodes)
}

func KargersWorker(id int, graph map[string][]string, weights map[string]map[string]int, nodes []string, max_runs int, results chan<- [2]int, wg *sync.WaitGroup, do_runs int) {
    for i := 0; i < do_runs; i++ {
        // Reducde the size of the graph by contracting random edges
        new_graph, new_weights, new_nodes := RandomlyContractEdges(graph, weights, nodes, 4)
        min_cut, min_groups_size := KargersRun(new_graph, new_weights, new_nodes)
        results <- [2]int{min_cut, min_groups_size}
    }
    wg.Done()
}

func KargersWorkerAlt(id int, graph map[string][]string, weights map[string]map[string]int, nodes []string, max_runs int, results chan<- [2]int, wg *sync.WaitGroup, do_runs int) {
    for i := 0; i < do_runs; i++ {
        // Reducde the size of the graph by contracting random edges
        min_cut, min_groups_size := KargersRunAlt(graph, weights, nodes)
        results <- [2]int{min_cut, min_groups_size}
    }
    wg.Done()
}

func RandomlyContractEdges(graph map[string][]string, weights map[string]map[string]int, nodes []string, max_contracts int) (map[string][]string, map[string]map[string]int, []string) {
    // Copy nodes 
    run_nodes := []string{}
    for _, node := range nodes {
        run_nodes = append(run_nodes, node)
    }

	// Shuffle the slice
	rand.Shuffle(len(run_nodes), func(i, j int) {
	    run_nodes[i], run_nodes[j] = run_nodes[j], run_nodes[i]
	})

    merges := [][]string{}
    for len(run_nodes) > 2 {
        nbr := rand.Intn(max_contracts) + 1
        if nbr < len(run_nodes) {
            to_merge := []string{}
            for i := 0; i < nbr; i++ {
                to_merge = append(to_merge, run_nodes[i])
            }
            merges = append(merges, to_merge)
            run_nodes = run_nodes[nbr:]
        } else {
            merges = append(merges, run_nodes)
            run_nodes = []string{}
        }
    }
    new_graph := make(map[string][]string)
    new_weights := make(map[string]map[string]int)
    new_nodes := []string{}
    for _, node := range nodes {
        new_graph[node] = []string{}
        new_weights[node] = make(map[string]int)
        new_nodes = append(new_nodes, node)
        for _, other_node := range nodes {
            new_weights[node][other_node] = 0
        }
        for _, neighbour := range graph[node] {
            new_graph[node] = append(new_graph[node], neighbour)
            new_weights[node][neighbour] = weights[node][neighbour]
        }
    }
    for _, merge := range merges {
        if len(merge) > 1 {
            new_graph, new_weights, new_nodes = mergeNodesGroup(new_graph, new_weights, new_nodes, merge)
        }
    }
    return new_graph, new_weights, new_nodes
}


func KargersMinCut(graph map[string][]string, weights map[string]map[string]int, nodes []string, max_runs int, numWorkers int) (int, int) {
    
    // Create channels
    var wg sync.WaitGroup
    results := make(chan [2]int, numWorkers)
    final_results := make(chan [2]int)

    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go KargersWorkerAlt(i, graph, weights, nodes, max_runs, results, &wg, 1 + max_runs / numWorkers)
    }

    // Collect results in a separate goroutine
    go func() {
        wg.Wait()
        close(results)
    }()

    // Aggregate results
    go func() {
        min_cut := math.MaxInt
        min_groups_size := 0

        for result := range results {
            if result[0] < min_cut {
                min_cut = result[0]
                min_groups_size = result[1]
            }
        }
        final_results <- [2]int{min_cut, min_groups_size}
        close(final_results)
    }()

    // Read from the final results
    result := <-final_results
    return result[0], result[1]
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
    graph := make(map[string][]string)
    nodes := []string{}
    connections := [][2]string{}
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
            // Split the line into the two parts
            parts := strings.Split(cleaned_line, ": ")
            
            left := parts[0]
            right := strings.TrimSpace(parts[1])
            // Split the right part into the individual elements
            right_parts := strings.Split(right, " ")
            graph[left] = []string{}
            nodes = append(nodes, left)
            for _, part := range right_parts {
                graph[left] = append(graph[left], part)
                nodes = append(nodes, part)
                connections = append(connections, [2]string{left, part})
            }
		}
	}
    nodes = RemoveDuplicates(nodes)
    for k, v := range graph {
        for _, val := range v {
            graph[val] = append(graph[val], k)
        }
    }
    fmt.Println("Nodes", len(nodes))
    fmt.Println("Connections", len(connections))

    // Create workers
    num_workers := cores

    weights := make(map[string]map[string]int)
    for _, node := range nodes {
        weights[node] = make(map[string]int)
        for _, other_node := range nodes {
            weights[node][other_node] = 0
        }
        for _, neighbour := range graph[node] {
            weights[node][neighbour] = 1
        }
    }

    runs := 0
    max_runs := 100
    min_cut, min_groups_size := KargersMinCut(graph, weights, nodes, max_runs, num_workers)
    runs++

    fmt.Println("Runs", runs * max_runs)
    fmt.Println("Min cut", min_cut)
    fmt.Println("Min groups size", min_groups_size)
    fmt.Println("Running until min cut is 3")
    fmt.Println("")
    for min_cut != 3 {
        min_cut, min_groups_size = KargersMinCut(graph, weights, nodes, max_runs, num_workers)
        runs++
        fmt.Println("Runs", runs * max_runs)
        fmt.Println("Min cut", min_cut)
        fmt.Println("Min groups size", min_groups_size)
        fmt.Println("")
    }
}
