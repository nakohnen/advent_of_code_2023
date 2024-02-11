package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ToInt(in string) int {
    res, err := strconv.Atoi(in)
    if err != nil {
        fmt.Println("Could not convert to int:", in)
        os.Exit(1)
    }
    return res
}

func ToRune(val int) rune {
	return rune('0' + val)
}

func AllElementsSame[T comparable](slice []T) bool {
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			return false
		}
	}
	return true
}

// Function to find prime factors of a number
func PrimeFactors(n int) map[int]int {
	factors := make(map[int]int)
	// Count the number of 2s that divide n
	for n%2 == 0 {
		factors[2]++
		n = n / 2
	}
	// n must be odd at this point. So start from 3 and iterate until sqrt(n)
	for i := 3; i <= int(math.Sqrt(float64(n))); i = i + 2 {
		// While i divides n, count i and divide n
		for n%i == 0 {
			factors[i]++
			n = n / i
		}
	}
	// If n is a prime number greater than 2
	if n > 2 {
		factors[n]++
	}
	return factors
}

// Function to find LCM of an array of integers
func FindLCM(arr []int) int {
	overallFactors := make(map[int]int)
	for _, num := range arr {
		// Get prime factors of each number
		primeFactorsOfNum := PrimeFactors(num)
		for prime, power := range primeFactorsOfNum {
			if currentPower, exists := overallFactors[prime]; !exists || power > currentPower {
				// Store the highest power of each prime
				overallFactors[prime] = power
			}
		}
	}
	// Calculate LCM by multiplying the highest powers of all primes
	lcm := 1
	for prime, power := range overallFactors {
		lcm *= int(math.Pow(float64(prime), float64(power)))
	}
	return lcm
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func HasPrefix(line string, sub string) bool {
	if len(line) < len(sub) {
		return false
	}
	return line[0:len(sub)] == sub
}

func ReplaceCharacters(line string, rem string, rep string) string {
	result := line
	for _, c := range rem {
		result = strings.ReplaceAll(result, string(c), rep)
	}
	return result
}

func SumIntSlice(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}
	return r
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

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, i := range slice {
		if !seen[i] {
			result = append(result, i)
			seen[i] = true
		}
	}
	return result
}

func Factorial (n int) int {
    if n == 0 {
        return 1
    }
    res := n * Factorial(n-1)
    if res < 0 {
        fmt.Println("Overflow")
        os.Exit(1)
    }
    return res
}
