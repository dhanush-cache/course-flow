package mosh

import "time"

type CourseResponse struct {
	PageProps CoursePageProps `json:"pageProps"`
}
type CoursesResponse struct {
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

type Token struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}
