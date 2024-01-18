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
// type Task struct {
// 	Id          int    `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	Status      string `json:"status"`
// }

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

	var taskapi TaskApi
	taskapi.router = gin.Default()
	go taskapi.serve()
	apiUrl = "http://0.0.0.0:9090/"
	r := gin.Default()
	templates, te = template.ParseGlob("templates/*.html")
	if te != nil {
		panic(te)
	}
	// r.SetHTMLTemplate(tmpl)
	tg := r.Group("/tasks")
	r.GET("/", index)
	r.GET("/newtask", taskForm)

	r.POST("/edit/:id", editTaskFromForm)
	r.GET("/edit/:id", editForm)

	tg.POST("/", newTaskFromForm)
	tg.POST("/delete/:id", deleteTask)
	tg.POST("/statuschange/:id", statusChange)

	r.Run(":8000")

}

func deleteTask(c *gin.Context) {
	durl := make(chan string)
	go func() {
		durl <- fmt.Sprintf(apiUrl+"tasks/%s", c.Param("id"))
	}()
	req, err := http.NewRequest("DELETE", <-durl, nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{}
	_, e := client.Do(req)
	if e != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusFound, "/")
}

func taskForm(c *gin.Context) {

	c.File("static/createtask.html")
}

func statusChange(c *gin.Context) {
	usurl := make(chan string, 1)
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	t := returnTaskAsJSON(c, id)
	defer close(usurl)
	go func(t Task) {

		fmt.Println(t)
		var newStatus string
		if t.Status == "incomplete" {
			newStatus = "complete"
		} else {
			newStatus = "incomplete"
		}
		usurl <- fmt.Sprintf(apiUrl+"tasks/?status=%s&id=%d", newStatus, id)
	}(t)
	statusurl := <-usurl
	_, err := http.PostForm(statusurl, url.Values{})
	if err != nil {
		panic(err)
	}
	c.Redirect(http.StatusFound, "/")
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
	c.Redirect(http.StatusFound, "/")
}

func editForm(c *gin.Context) {
	t := make(chan Task)
	defer close(t)
	go func() {
		id, _ := strconv.Atoi(c.Param("id"))
		t <- returnTaskAsJSON(c, id)
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

		turl <- fmt.Sprintf(apiUrl+"tasks/?id=%s&title=%s&description=%s&status=%s", url.QueryEscape(c.Param("id")), url.QueryEscape(c.PostForm("title")), url.QueryEscape(c.PostForm("description")), url.QueryEscape(c.PostForm("status")))
		// url.QueryEscape(id, c.PostForm("title")), url.QueryEscape(c.PostForm("description")))
	}()
	_, err := http.PostForm(<-turl, url.Values{})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	c.Redirect(http.StatusFound, "/")
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

func returnTaskAsJSON(c *gin.Context, id int) Task {
	fmt.Println("return task")
	fmt.Println(id)
	resp, err := http.Get(fmt.Sprintf(apiUrl+"tasks/%d", id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	defer resp.Body.Close()

	// Check if the response status code is successful (2xx)
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("API request failed with status code %d", resp.StatusCode)})

	}
	var t Task
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	return t

}
