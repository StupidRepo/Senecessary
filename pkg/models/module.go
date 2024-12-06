package models

import "time"

type Module struct {
	ModuleId  string `json:"moduleId"`
	ContentId string `json:"contentId"`
	CourseId  string `json:"courseId"`
	SectionId string `json:"sectionId"`
	SessionId string `json:"sessionId"`

	ModuleOrder int    `json:"moduleOrder"`
	ModuleType  string `json:"moduleType"`

	Completed bool `json:"completed"`
	GaveUp    bool `json:"gaveUp"`

	Content struct{} `json:"content"`

	Score     int  `json:"score"`
	Submitted bool `json:"submitted"`

	TestingActive bool `json:"testingActive"`

	TimeFinished time.Time `json:"timeFinished"`
	TimeStarted  time.Time `json:"timeStarted"`
}
