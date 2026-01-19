package dreamsofcode

type Course struct {
	Name     string    `json:"name"`
	Sections []Section `json:"curriculum"`
}

type Section struct {
	Name    string   `json:"name"`
	Lessons []Lesson `json:"lessons"`
}

type Lesson struct {
	Name string `json:"name"`
}
