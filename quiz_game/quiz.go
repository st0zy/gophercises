package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var FileNotFoundError = errors.New("quiz file does not exist")
var IncorrectFormatOfQuiz = errors.New("format of the quiz paper is incorrect")

type Quiz struct {
	questions       []QuestionAnswer
	student         string
	currentScore    int
	timer           *time.Timer
	currentQuestion QuestionAnswer
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
	defer quizFile.Close()
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
	writer.Write([]byte("Press any key to start the quiz\n"))
	answerReader.ReadLine()
	answer := make(chan string)
	q.timer = time.NewTimer(time.Second * 5)

	for _, question := range q.questions {
		writer.Write([]byte(question.question + "\n"))
		q.currentQuestion = question

		go func() {
			answer <- q.readAnswer(answerReader)
		}()
		select {
		case ans := <-answer:
			if ans == q.currentQuestion.expectedAnswer {
				q.currentScore++
			}
		case <-q.timer.C:
			writer.Write([]byte("Time Up.\n"))
			q.endQuiz(writer)
			return
		}
	}
	q.endQuiz(writer)
}

func (q *Quiz) endQuiz(writer io.Writer) {
	writer.Write([]byte("The quiz has now ended\n"))
	writer.Write([]byte(fmt.Sprint(q.student, ", Your score is  ", q.currentScore, "\n")))
	writer.Write([]byte(fmt.Sprint("Thank you for attending the quiz,", q.student, ".\n")))
}

func (q *Quiz) readAnswer(answerReader *bufio.Reader) string {
	answerFromStudent, _, err := answerReader.ReadLine()
	if err != nil {
		fmt.Println(err)
	}
	return string(answerFromStudent)

}
