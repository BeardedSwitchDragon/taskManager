package main

import (

	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
)
func tasksToMaps(tasks []Task) map[string]string {
	m := make(map[string]string)

	for _, t := range tasks {
		//Converts to JSON format
		fmt.Println(t)
		m[strconv.Itoa(t.id)] = "{title:" + t.title + "," + "description:" + t.description + "," + "status:" + t.status + "}"

	}
	return m
}
func main() {
	r := gin.Default()
	r.GET("/fetch/all", func(c *gin.Context) {
		var f Filter
		f.unspecified()
		tasks := getTasks(f)
		// fmt.Println(yikes)
		c.JSON(200, tasksToMaps(tasks))
	})

	r.Run()
}