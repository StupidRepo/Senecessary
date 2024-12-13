package models

type CoursesSectionsResponse struct {
	Course   Course    `json:"course"`
	Sections []Section `json:"sections"`

	Count int `json:"count"`
}

type GetSignedCourseURLResponse struct {
	URL string `json:"url"`
}

type Course struct {
	ID       string `json:"id"`
	CourseID string `json:"courseId"`

	SectionIds []string `json:"sectionIds"`

	Name        string   `json:"name"`
	Description string   `json:"description"`
	Authors     []string `json:"authors"`

	StartingSectionId string `json:"startingSectionId"`
}
