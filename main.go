package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func main() {
	filename, timelimit, shuffle := getPointersOfFlagValues()
	records := getCsvContent(filename)
	var correct int64
	ans := make(chan int)
	fmt.Print("Press ENTER to start the quiz!")
	fmt.Scanln()
	go quiz(&correct, records, shuffle, ans)
	go timer(timelimit, &correct, ans)
	fmt.Println("Score:", <-ans, "/", len(*records))
}

func getPointersOfFlagValues() (*string, *int, *bool) {
	filename := flag.String("f", "problems.csv", "The filename of the csv with the quiz questions.")
	timelimit := flag.Int("t", 30, "The timelimit in seconds.")
	shuffle := flag.Bool("s", false, "Shuffle the order in which the questions are asked.")
	flag.Parse()
	return filename, timelimit, shuffle
}

func getCsvContent(filename *string) *[][]string {
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return &records
}

func quiz(correct *int64, records *[][]string, shuffle *bool, ans chan int) {
	if *shuffle {
		rand.Seed(time.Now().Unix())
		for _, i := range rand.Perm(len(*records)) {
			askQuestion(&(*records)[i], correct)
		}
	} else {
		for _, row := range *records {
			askQuestion(&row, correct)
		}
	}
	ans <- int(atomic.LoadInt64(correct))
}

func askQuestion(row *[]string, correct *int64) {
	fmt.Print((*row)[0], " ")
	var ans string
	fmt.Scanln(&ans)
	ans = strings.ToLower(strings.ReplaceAll(ans, " ", ""))
	if ans == (*row)[1] {
		atomic.AddInt64(correct, 1)
	}
}

func timer(seconds *int, correct *int64, ans chan int) {
	time.Sleep(time.Duration(*seconds) * time.Second)
	fmt.Println()
	ans <- int(atomic.LoadInt64(correct))
}
