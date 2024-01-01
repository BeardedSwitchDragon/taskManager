package main

import (
	"fmt"
)

func main() {

	testTask := Task{
		id: 0001,
		title: "erm this is a test!!!",
		description: "im testing this task (description)",
		status: "incomplete",

	}

	writeDB(testTask)
	fmt.Println("Hello world!")

}

