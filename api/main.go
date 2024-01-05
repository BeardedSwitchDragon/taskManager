package main

import (
	"encoding/json"
	// "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	fetch := r.Group("/fetch")
	// {
	// 	r.GET("/all", func(c *gin.Context) {
	// 		j := make(chan []byte)
	// 		// jpointer := &j
	// 		go func() {
	// 			defer close(j)

	// 			var f Filter
	// 			f.unspecified()
	// 			tasks := getTasks(f)

	// 			d, _ := json.Marshal(tasks)
	// 			d, _ = json.MarshalIndent(tasks, "", " ")
	// 			j <- d
	// 			fmt.Println(d)

	// 		}()
	// 		// fmt.Println(*jpointer)

	// 		c.JSON(200, string(<-j))

	// 	})

	// }
	fetch.GET("/",fetchTasksApi)

	//DELETE a certain task

	r.DELETE("/task/:id", func(c *gin.Context) {
		go func() {
			id, _ := strconv.Atoi(c.Param("id"))
			err := deleteTask(id)
			if err != nil {

				c.String(http.StatusBadRequest, "something went wrnog")
			} else {
				c.String(http.StatusOK, "successfully deleted")
			}
		}()

	})

	go func() {
		if err := r.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	// Keep the program running
	select {}
}

func fetchTasksApi(c *gin.Context) {
	j := make(chan []byte)
	go func() {
		defer close(j)
		f := newFilter()
		t, d, s := c.Query("title"), c.Query("description"), c.Query("status")

		if t != "<nil>" {
			f.Title = t
		}
		if d != "<nil>" {
			f.Description = d
		}
		if s != "<nil>" {
			f.Status = s
		}
		tasks := getTasks(f)
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
