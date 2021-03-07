package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question string
	answer string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv",
		"a cvs file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed  to open CSV file: %s\n", *csvFileName))

	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)

	correct := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s =\n", i+1, problem.question)
		var answer string
		_, _ = fmt.Scanf("%s\n", &answer)
		if answer == problem.answer {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
