package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

var FileNotFoundError = errors.New("quiz file does not exist")
var IncorrectFormatOfQuiz = errors.New("format of the quiz paper is incorrect")

type Quiz struct {
	questions    []QuestionAnswer
	student      string
	currentScore int
}

type QuestionAnswer struct {
	question       string
	expectedAnswer string
}

func NewQuizFromPath(path string, student string) (*Quiz, error) {
	quizFile, err := os.Open(path)
	if err != nil {
		return nil, FileNotFoundError
	}
	reader := csv.NewReader(quizFile)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, IncorrectFormatOfQuiz
	}

	var questions []QuestionAnswer
	for _, record := range records {
		questions = append(questions, QuestionAnswer{record[0], record[1]})
	}
	//Do we need to randomize the questions?
	return &Quiz{
		questions:    questions,
		currentScore: 0,
		student:      student,
	}, nil

}

func (q *Quiz) StartQuiz(writer io.Writer, reader io.Reader) {
	answerReader := bufio.NewReader(reader)
	writer.Write([]byte(fmt.Sprint("Welcome to the quiz,", q.student, "\n")))
	writer.Write([]byte("The quiz will start now....\n"))

	for _, question := range q.questions {
		writer.Write([]byte(question.question))
		answerFromStudent, _, err := answerReader.ReadLine()
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(answerFromStudent)
		if string(answerFromStudent) == question.expectedAnswer {
			q.currentScore++
		}
	}

	writer.Write([]byte("The quiz has now ended\n"))
	writer.Write([]byte(fmt.Sprint(q.student, ", Your score is  ", q.currentScore, "\n")))
	writer.Write([]byte(fmt.Sprint("Thank you for attending the quiz,", q.student, ".\n")))

}
