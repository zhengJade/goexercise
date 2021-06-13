package main

import (
	"encoding/csv"
	"os"
	"fmt"
	"strings"
	"strconv"
)

type Quiz struct {
	question string
	anwser string
}

func QuizReader(filePath string) []Quiz {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var quizs []Quiz
	for _, record := range records {
		quiz := Quiz{question: record[0], anwser: record[1]}
		quizs = append(quizs, quiz)
	}
	return quizs
}

func compute(quiz Quiz) bool {
	question := quiz.question
	exp := strings.Split(question, "+")
	sum := 0
	for _, num := range exp {
		num, _ := strconv.Atoi(num)
		sum += num
	}
	anwser, _ := strconv.Atoi(quiz.anwser)
	if sum == anwser {
		return true
	} else {
		return false
	}
}

func main() {
	filePath := "/Users/jade/projects/go-test/quiz/problems.csv"
	quizs := QuizReader(filePath)
	fmt.Println(len(quizs))
	total := 0
	for _, quiz := range quizs {
		if compute(quiz) {
			total++
		}
	}
	fmt.Println(total)
}