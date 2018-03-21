package main

import (
    "fmt"
)

func main() {
    array := []int{-7, 1, 5, 2, -4, 3, 0}
    fmt.Println("Answer:", solution(array, len(array)))
}

func solution(ints []int, n int) int{

    if n == 1 {
        return -1
    }

    if n == 0 {
        return 0
    }

    totalSum := total(ints)
    sum := 0

    for i:=0; i<n; i++ {
        totalSum -= ints[i]
        fmt.Printf("Total sum %d | current sum %d\n", totalSum, sum)
        if sum == totalSum {
            return i
        }

        sum += ints[i]
    }

    return -1
}

func total(ints []int) int {
    total := 0
    for _, i := range ints {
        total += i
    }

    return total
}
