package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	fileFlag := flag.String("file", "problems.csv", "path for quiz file")
	timerFlag := flag.Int("time", 5, "set time limit for quit")
	flag.Parse()

	f, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("failed to read from file: %v", err)
	}

	timer := time.NewTimer(time.Duration(*timerFlag) * time.Second)
	count := 0
	ans := ""

outerLoop:
	for i, record := range records {
		fmt.Printf("#%d: %s: ", i+1, record[0])
		ansChan := make(chan string)
		go func() {
			ans := ""
			fmt.Scan(&ans)
			ansChan <- ans
		}()

		select {
		case <-timer.C:
			fmt.Print("\n")
			break outerLoop
		case ans = <-ansChan:
			if ans == record[1] {
				count++
			}
		}
	}

	fmt.Println("Out of ", len(records), "questions you got ", count, "correct")
}
