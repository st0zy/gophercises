package main

import (
	"flag"
	"os"
)

var path string
var student string

func main() {
	// fmt.Println("Hello world.")

	flag.StringVar(&path, "path", "problems.csv", "Provide a new path for the set of questions")
	flag.StringVar(&student, "student", "Student1", "Provide the name of the student")
	flag.Parse()

	q, err := NewQuizFromPath(path, student)
	if err != nil {
		panic(err)
	}
	q.StartQuiz(os.Stdout, os.Stdin)

}
