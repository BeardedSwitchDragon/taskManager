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
	//Fetch group, task group, admin group respectively.
	fg := r.Group("/fetch")
	tg := r.Group("/task")
	ag := r.Group("/admin")

	ag.POST("/newdb", func(c *gin.Context) {
		go func() {
			createDB()
		}()
	})

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
	m := make(chan string)
	defer close(m)
	s := make(chan int)
	
	go func() {
		
		id, _ := strconv.Atoi(c.Param("id"))
		err := deleteTask(id)
		if err != nil {

			s <- http.StatusBadRequest
			m <- "something went wrong!"
		} else {
			s <- http.StatusOK
			m <- "successfully deleted!"
		}
	}()
	c.String(<-s, <-m)

}

func createTaskApi(c *gin.Context) {
	go func() {
		// title, d, s := c.Query("title"), c.Query("description"), c.Query("status")
		fmt.Println("hi")
		var id int
		var t Task
		var err error
		if c.Query("id") != "" {
			//Logic if a task already exists (updates existing one)
			
			id, _ = strconv.Atoi(c.Query("id"))
			t, err = getTask(id)

			if err != nil {
				c.String(http.StatusNotFound, "404 Not Found: The requested resource was not found.")
				return
			}

			if c.Query("title") != "" {
				t.Title = c.Query("title")
			}
			if c.Query("description") != "" {
				t.Description = c.Query("description")
			} 
			if c.Query("status") != "" {
				t.Status = c.Query("status")
			}
			c.String(http.StatusOK, "successfully updated")

		} else {
			//Logic if a new Task is to be created
			var empty Filter
			
			empty.unspecified()
			existingTasks := getTasks(empty)
			id = len(existingTasks) + 1
			t = Task {
				Id: id,
				Title: c.Query("title"),
				Description: c.Query("description"),
				Status: c.Query("status"),
			}
			writeDB(t)
			c.String(http.StatusOK, "successfully created task")
		}
		
		
		

	}()
	
}


//fetch task function
func fetchTasksApi(c *gin.Context) {
	j := make(chan []byte)
	go func() {
		defer close(j)
		f := newFilter()
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
		if (t == "") && (d == "") && (s == "") {
			f.unspecified()
			fmt.Println(f)
		}
		fmt.Println(f)
		tasks := getTasks(f)

		data, _ := json.Marshal(tasks)
		if len(tasks) == 0 {
			data = []byte("[]")
		}
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
