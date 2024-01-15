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

type Workflow struct {
	name   string
	rules  []WorkflowRule
	target string
}

func shortenWorkflow(workflows []Workflow) []Workflow {
	result := []Workflow{}
	for _, wf := range workflows {
		new_end_id := len(wf.rules)
		for i := len(wf.rules) - 1; i >= 0; i-- {
			if wf.rules[i].target != wf.target {
				new_end_id = i + 1
				break
			}
		}
		new_wf := Workflow{name: wf.name, rules: wf.rules[:new_end_id], target: wf.target}
		result = append(result, new_wf)

	}
	return result
}

func simplifyRules(workflows []Workflow) []Workflow {
	result_rules := []Workflow{}

	// Remove redundant rules
	redundant := []Workflow{}
	redundant_map := make(map[string]bool)
	redundant_target := make(map[string]string)
	for _, wf := range workflows {
		rules := wf.rules
		end := wf.target

		found := true
		for _, sub_rule := range rules {
			if sub_rule.target != end {
				found = false
				break
			}
		}
		if found {
			fmt.Printf("Removing rule %v\n", wf.name)
			redundant = append(redundant, wf)
			redundant_map[wf.name] = true
			redundant_target[wf.name] = wf.target
		}
	}

	for _, wf := range workflows {
		found := false
		for _, red := range redundant {
			if wf.name == red.name {
				found = true
				break
			}
		}
		if !found {
			name := wf.name
			rules := []WorkflowRule{}
			target := wf.target
			if redundant_map[wf.target] {
				target = redundant_target[wf.target]
			}
			for i := 0; i < len(wf.rules); i++ {
				rule := wf.rules[i]
				rule_target := wf.rules[i].target
				if redundant_map[rule_target] {
					rule_target = redundant_target[rule_target]
				}
				new_rule := rule
				new_rule.target = rule_target
				rules = append(rules, new_rule)
			}
			new_wf := Workflow{name: name, target: target, rules: rules}
			result_rules = append(result_rules, new_wf)
		}

	}
	return shortenWorkflow(result_rules)
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

func workWorkflow(name string, queue <-chan Gear, rules []WorkflowRule, end string, other map[string]chan Gear, wg *sync.WaitGroup, resultChan chan<- int) {
	result := 0
	for gear := range queue {
		target := evaluateRule(name, rules, end, gear)
		if target == "R" || target == "A" {
			if debug {
				fmt.Printf("Gear %v ended in the end pile %v.\n", gear, target)
			}
			wg.Done()
			if target == "A" {
				result++
			}
		} else {
			other[target] <- gear
		}
	}
	resultChan <- result
}

type GearRange struct {
	attr map[string][2]int
}

var gearAttr = [4]string{"a", "x", "s", "m"}

func splitGearRangeAlong(gr GearRange, attr string, value int, op string) ([]GearRange, int) {
	result := []GearRange{}
	sit := 0
	target_range := gr.attr[attr]
	if target_range[0] <= value && value <= target_range[1] {
		lower_range := GearRange{make(map[string][2]int)}
		upper_range := GearRange{make(map[string][2]int)}
		for _, a := range gearAttr {
			if a != attr {
				lower_range.attr[a] = gr.attr[a]
				upper_range.attr[a] = gr.attr[a]
			}
		}
		if op == "<" {
			lower_range.attr[attr] = [2]int{target_range[0], value - 1}
			upper_range.attr[attr] = [2]int{value, target_range[1]}
			result = append(result, lower_range)
			result = append(result, upper_range)
		} else if op == ">" {
			lower_range.attr[attr] = [2]int{target_range[0], value}
			upper_range.attr[attr] = [2]int{value + 1, target_range[1]}
			result = append(result, lower_range)
			result = append(result, upper_range)
		}

	} else {
		if target_range[0] > value {
			sit = 1
		} else if target_range[1] < value {
			sit = 2
		}
		result = append(result, gr)
	}
	// sit = 0 we split
	// sit = 1 our range is higher than value
	// sit = 2 our range is lower than value
	return result, sit
}

func evaluateRulesRanges(name string, workflows map[string]Workflow, g_range []GearRange) []GearRange {
	accepted := []GearRange{}
	rejected := []GearRange{}
	to_work := []GearRange{}
	for _, r := range g_range {
		to_work = append(to_work, r)
	}
	wf := workflows[name]
	for _, rule := range wf.rules {
		new_work := []GearRange{}
		for _, r := range to_work {
			candidates, sit := splitGearRangeAlong(r, rule.attr, rule.value, rule.op)
			switch sit {
			case 0:
				inside_id := 1
				outside_id := 0
				if rule.op == "<" {
					inside_id = 0
					outside_id = 1
				}
				switch rule.target {
				case "A":
					accepted = append(accepted, candidates[inside_id])
				case "R":
					rejected = append(rejected, candidates[inside_id])
				default:
					for _, new_r := range evaluateRulesRanges(rule.target, workflows, []GearRange{candidates[inside_id]}) {
						accepted = append(accepted, new_r)
					}
				}
				new_work = append(new_work, candidates[outside_id])
			case 1: // rule.val is higher than the range
				switch rule.op {
				case ">": // range > val
					switch rule.target {
					case "A":
						accepted = append(accepted, candidates[0])
					case "R":
						rejected = append(rejected, candidates[0])
					default:
						for _, new_r := range evaluateRulesRanges(rule.target, workflows, []GearRange{candidates[0]}) {
							accepted = append(accepted, new_r)
						}
					}
				case "<":
					new_work = append(new_work, candidates[0])
				}
			case 2: // rule.val is lower than the range
				switch rule.op {
				case "<": // range < val
					switch rule.target {
					case "A":
						accepted = append(accepted, candidates[0])
					case "R":
						rejected = append(rejected, candidates[0])
					default:
						for _, new_r := range evaluateRulesRanges(rule.target, workflows, []GearRange{candidates[0]}) {
							accepted = append(accepted, new_r)
						}
					}
				case ">":
					new_work = append(new_work, candidates[0])
				}
			}
		}
		to_work = new_work
	}
	switch wf.target {
	case "R":
		for _, r := range to_work {
			rejected = append(rejected, r)
		}
	case "A":
		for _, r := range to_work {
			accepted = append(accepted, r)
		}
	default:
		for _, r := range evaluateRulesRanges(wf.target, workflows, to_work) {
			accepted = append(accepted, r)
		}
	}
	return accepted
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
	workflows := []Workflow{}
	rules_part := true
	workflows_map := make(map[string]Workflow)
	for _, line := range strings.Split(string(dat), "\n") {
		cleaned_line := strings.TrimSpace(line)
		if len(cleaned_line) > 0 {
			if debug {
				fmt.Printf("%v\n", cleaned_line)
			}
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
				wf := Workflow{name: name, target: end, rules: rules}
				workflows = append(workflows, wf)
				workflows_map[wf.name] = wf
			}
		} else {
			rules_part = false
		}
	}

	for _, wf := range workflows {
		fmt.Printf("%v: %v %v\n", wf.name, wf.rules, wf.target)
	}
	fmt.Println("")
	workflows = simplifyRules(workflows)
	for _, wf := range workflows {
		fmt.Printf("%v: %v %v\n", wf.name, wf.rules, wf.target)
		workflows_map[wf.name] = wf
	}

	start_range := GearRange{attr: make(map[string][2]int)}
	for _, attr := range gearAttr {
		start_range.attr[attr] = [2]int{1, 4000}
	}
	result := 0
	for _, r := range evaluateRulesRanges("in", workflows_map, []GearRange{start_range}) {
		my_val := 1
		for _, attr := range gearAttr {
			lower := r.attr[attr][0]
			higher := r.attr[attr][1]
			my_val *= higher - lower + 1

		}
		result += my_val
	}

	fmt.Println("Final results:", result)

	fmt.Println("")
}
