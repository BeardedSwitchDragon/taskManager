package main
import "fmt"
type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
}

type Filter struct {
	Title string
	Description string
	Status string
}

func (f *Filter) unspecified() {
	f.Title, f.Description, f.Status = "", "", ""
	fmt.Println(f.Title)
}


//Uses undefined character as default that is universally understood to be an invalid character type.
func newFilter() Filter {
	return Filter{
		Title: "𑨩",
		Description: "𑨩",
		Status: "𑨩",
	}
} 