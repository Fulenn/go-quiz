package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Question struct {
	QuestionText string
	Answer       string
	Choices      []string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the name of the quiz file (e.g., quiz.csv): ")
	scanner.Scan()
	filename := scanner.Text()

	fmt.Print("Enter the amount of time per question (in seconds): ")
	scanner.Scan()
	timeLimitStr := scanner.Text()

	timeLimit, err := time.ParseDuration(timeLimitStr + "s")
	if err != nil {
		fmt.Println("Invalid time duration. Exiting.")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var questions []Question
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		questions = append(questions, Question{
			QuestionText: record[0],
			Answer:       record[1],
			Choices:      record[2:],
		})
	}

	correctCount := 0
	timer := time.NewTimer(timeLimit)

	for i, q := range questions {
		fmt.Printf("Question %d: %s", i+1, q.QuestionText)
		for j, choice := range q.Choices {
			fmt.Printf("%d. %s ", j+1, choice)
		}
		fmt.Println()

		answerCh := make(chan string)
		go func() {
			scanner.Scan()
			answerCh <- scanner.Text()
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			fmt.Printf("You got %d out of %d questions correct.\n", correctCount, len(questions))
			return
		case answer := <-answerCh:
			selectedChoice, _ := strconv.Atoi(answer)
			if q.Choices[selectedChoice-1] == q.Answer {
				correctCount++
			}
		}
	}
	fmt.Printf("You got %d out of %d questions correct.\n", correctCount, len(questions))
}
