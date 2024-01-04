package main

type Task struct {
	id int
	title string
	description string
	status string
}

type Filter struct {
	title string
	description string
	status string
}

func (f Filter) unspecified() {
	f.title, f.description, f.status = "", "", ""
}


//Uses undefined character as default that is universally understood to be an invalid character type.
func newFilter() Filter {
	return Filter{
		title: "𑨩",
		description: "𑨩",
		status: "𑨩",
	}
} 