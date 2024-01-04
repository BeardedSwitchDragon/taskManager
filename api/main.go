package main

import (

	"github.com/gin-gonic/gin"
	// "fmt"
	"encoding/json"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	r.GET("/fetch/all", func(c *gin.Context) {
		var f Filter
		f.unspecified()
		tasks := getTasks(f)
		j, _ := json.Marshal(tasks)
		j, _ = json.MarshalIndent(tasks, "", " ")
		c.JSON(200, string(j))

	})


	//DELETE a certain task

	r.DELETE("/task/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := deleteTask(id)
		if err != nil {
	
			c.String(http.StatusBadRequest, "something went wrnog")
		} else {
			c.String(http.StatusOK, "successfully deleted")
		}
	})

	r.Run()
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