package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"time"
)

func reportResults(correct, total int) {
	fmt.Printf("\nYou scored %d/%d!\n", correct, total)
}

// deliver the quiz and collect results
func quiz(qs [][]string, correct *int, c chan int) {
	var ans string
	var answered int
	for _, q := range qs {
		fmt.Print(q[0] +  "=")
		fmt.Scan(&ans)
		if ans == q[1] {
			*correct = *correct + 1
		}
		c <- answered
	}
}

func main() {
	// open the csv file
	probs, _ := os.Open("problems.csv")

	// create a reader from csv package
	r := csv.NewReader(probs)
	r.FieldsPerRecord = 2

	// parse csv file
	qs, _ := r.ReadAll()

	var total int = len(qs)
	var correct int

	// create channel and timer for quiz
	c := make(chan int, 2)
	timer := time.NewTimer(3 * time.Second)
	
	// start quiz in background 
	// not sure if this should be a closure (check idioms)
	go quiz(qs, &correct, c)

	// select blocks on channels until first receive
	for {
		select {
		case <-timer.C:
			reportResults(correct, total)
			return
		case answered := <-c:
			if answered == total {
				reportResults(correct, total)
				return
			}
		}
	}
}
