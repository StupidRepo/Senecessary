package models

import "time"

type Status string

const (
	Complete   Status = "COMPLETE"
	Incomplete Status = "INCOMPLETE"
)

type AssignmentResponse struct {
	Items []Assignment `json:"items"`
	Count int          `json:"count"`
}

type Assignment struct {
	Id      string `json:"id"`
	ClassId string `json:"classId"`
	UserId  string `json:"userId"`

	NumAssignees int `json:"numAssignees"`

	TimeCreated time.Time `json:"timeCreated"`
	TimeUpdated time.Time `json:"timeUpdated"`

	StartDate time.Time `json:"startDate"`
	DueDate   time.Time `json:"dueDate"`
	Archived  bool      `json:"archived"`

	Name string `json:"name"`

	Spec AssigmentSpec `json:"spec"`

	Status Status `json:"status"`
}

type AssigmentSpec struct {
	CourseId    string   `json:"courseId"`
	SectionIds  []string `json:"sectionIds"`
	QuestionIds []string `json:"questionId"`
}
