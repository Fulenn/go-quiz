package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Question struct {
	QuestionText string
	Answer       string
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
		fmt.Println("Invalid time format. Please enter a valid duration (e.g., '30' for 30 seconds).")
		return
	}

	questionList := readFile(filename)
	quiz(questionList, timeLimit)
}

func readFile(s string) []Question {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	questionList := []Question{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // Exit the loop at the end of the file
		} else if err != nil {
			fmt.Println("Error:", err)
			panic(err)
		}

		question := Question{
			QuestionText: record[0],
			Answer:       record[1],
		}

		questionList = append(questionList, question)
	}

	return questionList
}

func quiz(q []Question, timeLimit time.Duration) {
	var counter int = 0
	scanner := bufio.NewScanner(os.Stdin)

	for _, question := range q {
		color.Cyan(question.QuestionText)

		// Create a timer with a 10-second duration
		timer := time.NewTimer(timeLimit)

		answerCh := make(chan string)

		// Use a goroutine to capture the user's input
		go func() {
			if scanner.Scan() {
				answerCh <- scanner.Text()
			}
		}()

		// Wait for either the user's input or the timer to expire
		select {
		case answer := <-answerCh:
			// The user answered in time
			if checkAnswer(answer, question.Answer) {
				counter++
			}
		case <-timer.C:
			// Time's up
			fmt.Println("Time's up! The correct answer was:", question.Answer)
		}
	}
	fmt.Printf("\nYou answered %d out of %d questions correctly.\n", counter, len(q))
}

func checkAnswer(answer string, correctAnswer string) bool {
	// Convert both the user's answer and the correct answer to lowercase for case-insensitive matching
	answer = strings.ToLower(answer)
	correctAnswer = strings.ToLower(correctAnswer)

	// Remove leading and trailing spaces from the user's answer
	answer = strings.TrimSpace(answer)

	// You can add more flexibility by using string similarity metrics like Levenshtein distance
	// or allowing synonyms, but here's a basic case-insensitive check
	return answer == correctAnswer
}
