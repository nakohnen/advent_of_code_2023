package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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

func toRune(val int) rune {
	return rune('0' + val)
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

func removeDuplicates(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, i := range slice {
		if !seen[i] {
			result = append(result, i)
			seen[i] = true
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

type Gear struct {
	x, m, a, s int
}

type WorkflowRule struct {
	attr, op string
	value    int
	target   string
}

func evaluateRule(name string, rules []WorkflowRule, end string, g Gear) string {
	for _, rule := range rules {
        found := false
		switch attr := rule.attr; attr {
		case "a":
			if rule.op == ">" {
				if g.a > rule.value {
					found = true
				}
			} else {
				if g.a < rule.value {
					found = true
				}
			}
		case "m":
			if rule.op == ">" {
				if g.m > rule.value {
					found = true
				}
			} else {
				if g.m < rule.value {
					found = true
				}
			}
		case "x":
			if rule.op == ">" {
				if g.x > rule.value {
					found = true
				}
			} else {
				if g.x < rule.value {
					found = true
				}
			}
		case "s":
			if rule.op == ">" {
				if g.s > rule.value {
					found = true
				}
			} else {
				if g.s < rule.value {
					found = true
				}
			}
		}
        if found {
            if debug {
                fmt.Printf("%v: Rule %v catched %v gear.\n", name, rule, g)
            }
            return rule.target
        }
	}
    if debug {
        fmt.Printf("%v: No rule %v catched %v gear.\n", name, rules, g)
    }
	return end
}

func workWorkflow(name string, queue <-chan Gear, rules []WorkflowRule, end string, other map[string]chan Gear, wg *sync.WaitGroup) {
	for gear := range queue {
		target := evaluateRule(name, rules, end, gear)
		other[target] <- gear
        if target == "R" || target == "A" {
            if debug {
                fmt.Printf("Gear %v ended in the end pile %v.\n", gear, target)
            }
            wg.Done()
        }
	}
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
	check(err)

	// Read file
	workflows := [][]WorkflowRule{}
	names := []string{}
	end_parts := []string{}
	gears := []Gear{}
	rules_part := true
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			fmt.Printf("%v\n", cleaned_line)
			if rules_part {
				rules := []WorkflowRule{}
				first_split := strings.Split(cleaned_line, "{")
				name := first_split[0]
				second_split := strings.Split(replaceCharacters(first_split[1], "}", ""), ",")
				end := second_split[len(second_split)-1]
				for _, raw_rule := range second_split[:len(second_split)-1] {
					third_split := strings.Split(raw_rule, ":")
					attr := third_split[0][0]
					op := third_split[0][1]
					value := toInt(third_split[0][2:])
					target := third_split[1]
					rules = append(rules, WorkflowRule{string(attr), string(op), value, target})
				}
				names = append(names, name)
				end_parts = append(end_parts, end)
				workflows = append(workflows, rules)
			} else {
				cleaned_line = replaceCharacters(cleaned_line, "{}", "")
				a, m, x, s := 0, 0, 0, 0
				for _, attr_raw := range strings.Split(cleaned_line, ",") {
					attr := string(attr_raw[0])
					value := toInt(attr_raw[2:])
					switch attr {
					case "a":
						a = value
					case "m":
						m = value
					case "x":
						x = value
					case "s":
						s = value
					}
				}
				gears = append(gears, Gear{a: a, m: m, x: x, s: s})
			}
		} else {
			rules_part = false
		}
	}

	var wg sync.WaitGroup
	work_queues := make(map[string]chan Gear)
	work_queues["R"] = make(chan Gear, len(gears))
	work_queues["A"] = make(chan Gear, len(gears))

	for i, name := range names {
		queue := make(chan Gear)
		work_queues[name] = queue
		go workWorkflow(name, queue, workflows[i], end_parts[i], work_queues, &wg)
	}

	for _, gear := range gears {
		wg.Add(1)
		work_queues["in"] <- gear
	}
    close(work_queues["in"])

	wg.Wait()
    fmt.Println("Wait over, closing channels")
	for _, name := range names {
        if name != "in" {
		    close(work_queues[name])
        }
	}
	close(work_queues["R"])
    close(work_queues["A"])

	// Collect results
	results := 0
	for gear := range work_queues["A"] {
		results += gear.a
		results += gear.m
		results += gear.s
		results += gear.x
	}

	fmt.Println("Final results:", results)

	fmt.Println("")
}
