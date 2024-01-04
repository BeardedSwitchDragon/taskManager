package main

import (
	"fmt"
)

func main() {

	// testTask := Task{
	// 	id: 0003,
	// 	title: "!!!",
	// 	description: "description)",
	// 	status: "incomplete",

	// }

	filter := newFilter()
	
	filter.title = "test"
	fmt.Println(filter)
	tasks := getTasks(filter)
	fmt.Println(tasks)

}

