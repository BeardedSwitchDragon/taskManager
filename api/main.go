package main

import (

	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	r.GET("/fetch/all", func(c *gin.Context) {
		j := make(chan []byte)
		// jpointer := &j
		go func() {
			defer close(j)
			
			var f Filter
			f.unspecified()
			tasks := getTasks(f)
			
			d, _ := json.Marshal(tasks)
			d, _ = json.MarshalIndent(tasks, "", " ")
			j <- d
			fmt.Println(d)
			
	
		}()
		// fmt.Println(*jpointer)
		
		c.JSON(200, string(<-j))
		
	})


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

// func tasksToMaps(tasks []Task) map[string]string {
// 	m := make(map[string]string)

// 	for _, t := range tasks {
// 		//Converts to JSON format
// 		fmt.Println(t)
// 		m[strconv.Itoa(t.id)] = "{title:" + t.title + "," + "description:" + t.description + "," + "status:" + t.status + "}"

// 	}
// 	return m
// }