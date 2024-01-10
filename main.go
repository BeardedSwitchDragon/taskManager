package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// redefining the Task type
type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type PageData struct {
	Tasks []Task

}

// Two seperate instances are being run - the api and the webpage. Both are being run locally on the same machine.
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
	r.GET("/newtask", taskForm)
	r.POST("/tasks", newTaskFromForm)
	r.POST("/edit/:id", editTaskFromForm)
	r.GET("/edit/:id", editForm)
	

	r.Run(":8000")

}

func taskForm(c *gin.Context) {

	c.File("static/createtask.html")
}

func newTaskFromForm(c *gin.Context) {
	turl := make(chan string)
	defer close(turl)
	go func() {

		turl <- fmt.Sprintf(apiUrl+"tasks/?title=%s&description=%s&status=incomplete",
			url.QueryEscape(c.PostForm("title")), url.QueryEscape(c.PostForm("description")))

	}()
	resp, err := http.PostForm(<-turl, url.Values{})
	if err != nil {
		panic(err)
	}

	body := make([]byte, 0)
	_, _ = resp.Body.Read(body)
	fmt.Println(resp.Status, string(body))
}

func editForm(c *gin.Context) {
	t := make(chan Task)
	defer close(t)
	go func() {
		id, _ := strconv.Atoi(c.Query("id"))
		t <- returnTasksAsJSON(c)[id]
	}()
	err := templates.ExecuteTemplate(c.Writer, "edittaskform.html", <-t)
	if err != nil {
		panic(err)
	}
}

func editTaskFromForm(c *gin.Context) {
	turl := make(chan string)
	defer close(turl)
	go func() {
		id, _ := strconv.Atoi(c.Query("id"))
		turl <-  fmt.Sprintf(apiUrl+"tasks/?id=%d&title=%s&description=%s&status=incomplete",id, url.QueryEscape(c.PostForm("title")), url.QueryEscape(c.PostForm("description")))
		// url.QueryEscape(id, c.PostForm("title")), url.QueryEscape(c.PostForm("description")))
	}()
	_, err := http.PostForm(<-turl, url.Values{})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}


func index(c *gin.Context) {
	tasks := make(chan []Task)
	defer close(tasks)
	tempTasks := returnTasksAsJSON(c)
	go func(ts []Task) {
		tasks <- ts
	}(tempTasks)
	err := templates.ExecuteTemplate(c.Writer, "index.html", PageData{Tasks: <-tasks})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	// c.HTML(http.StatusOK, "index", gin.H{"Tasks": <-tasks})

}


func returnTasksAsJSON(c *gin.Context) []Task {
	resp, err := http.Get(apiUrl + "tasks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		
	}
	defer resp.Body.Close()

	// Check if the response status code is successful (2xx)
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("API request failed with status code %d", resp.StatusCode)})
		
	}

	var tempTasks []Task
	err = json.NewDecoder(resp.Body).Decode(&tempTasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		
	}
	fmt.Println(tempTasks)

	return tempTasks

}