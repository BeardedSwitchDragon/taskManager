package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

//redefining the Task type
type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
}

type PageData struct {
	
	Tasks []Task
}

//Two seperate instances are being run - the api and the webpage. Both are being run locally on the same machine.
var apiUrl string

var templates *template.Template
var te error

func main() {

	// testTask := Task{
	// 	id: 0003,
	// 	title: "!!!",
	// 	description: "description)",
	// 	status: "incomplete",

	// }
	apiUrl = "http://0.0.0.0:8080/"
	r := gin.Default()
	templates, te = template.ParseGlob("templates/*.html")
	if te != nil {
		panic(te)
	}
	// r.SetHTMLTemplate(tmpl)
	r.GET("/", index)
	
	r.Run(":8000")

}

func index(c *gin.Context) {
	tasks := make(chan []Task)
	defer close(tasks)
	go func() {
		resp, err := http.Get(apiUrl + "tasks")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Check if the response status code is successful (2xx)
		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("API request failed with status code %d", resp.StatusCode)})
			return
		}

		var tempTasks []Task
		err = json.NewDecoder(resp.Body).Decode(&tempTasks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(tempTasks)

		tasks <- tempTasks

	}()
	err := templates.ExecuteTemplate(c.Writer, "index.html", PageData{Tasks: <-tasks})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	// c.HTML(http.StatusOK, "index", gin.H{"Tasks": <-tasks})

	
}