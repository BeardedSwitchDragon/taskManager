package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	//Fetch group, task group respectively.
	fg := r.Group("/fetch")
	tg := r.Group("/task")

	fg.GET("/", fetchTasksApi)

	//DELETE a certain task
	tg.DELETE("/:id", deleteTaskApi)
	go func() {
		if err := r.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	tg.POST("/", createTaskApi)

	select {}

}


//deletes task function
func deleteTaskApi(c *gin.Context) {
	go func() {
		id, _ := strconv.Atoi(c.Param("id"))
		err := deleteTask(id)
		if err != nil {

			c.String(http.StatusBadRequest, "something went wrnog")
		} else {
			c.String(http.StatusOK, "successfully deleted")
		}
	}()

}

func createTaskApi(c *gin.Context) {
	go func() {
		// title, d, s := c.Query("title"), c.Query("description"), c.Query("status")
		var empty Filter
		empty.unspecified()
		existingTasks := getTasks(empty)
		t := Task {
			Id: len(existingTasks) + 1,
			Title: c.Query("title"),
			Description: c.Query("description"),
			Status: c.Query("status"),
		}
		writeDB(t)

	}()
	
}


//fetch task function
func fetchTasksApi(c *gin.Context) {
	j := make(chan []byte)
	go func() {
		defer close(j)
		var f Filter
		t, d, s := c.Query("title"), c.Query("description"), c.Query("status")

		if t != "" {
			f.Title = t
		}
		if d != "" {
			f.Description = d
		}
		if s != "" {
			f.Status = s
		}
		fmt.Println(f)
		tasks := getTasks(f)
		fmt.Println(tasks)
		data, _ := json.Marshal(tasks)
		j <- data
	}()

	c.Data(200, "application/json; charset=utf-8", <-j)
}

// func tasksToMaps(tasks []Task) map[string]string {
// 	m := make(map[string]string)

// 	for _, t := range tasks {
// 		//Converts to JSON format
// 		fmt.Println(t)
// 		m[strconv.Itoa(t.id)] = "{title:" + t.title + "," + "description:" + t.description + "," + "status:" + t.status + "}"

// 	}
// 	return m
// }
