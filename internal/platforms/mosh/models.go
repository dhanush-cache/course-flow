package mosh

type CourseResponse struct {
	Props CourseProps `json:"props"`
}

type CourseProps struct {
	PageProps CoursePageProps `json:"pageProps"`
}

type CoursesResponse struct {
	Props CoursesProps `json:"props"`
}

type CoursesProps struct {
	PageProps CoursesPageProps `json:"pageProps"`
}

type CoursePageProps struct {
	Course Course `json:"course"`
}
type CoursesPageProps struct {
	Courses []Course `json:"courses"`
}

type Course struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Type           string    `json:"type"`
	BundleContents []int     `json:"bundleContents"`
	Curriculum     []Section `json:"curriculum"`
	Courses        []*Course `json:"courses"`
}

type Section struct {
	Name    string   `json:"name"`
	Lessons []Lesson `json:"lessons"`
}

type Lesson struct {
	Name    string `json:"name"`
	IsVideo int    `json:"type"`
}
