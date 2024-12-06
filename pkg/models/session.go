package models

import "time"

type SessionRequest struct {
	ClientVersion string `json:"clientVersion"`
	Platform      string `json:"platform"`

	Modules []Module `json:"modules"`
	Session Session  `json:"session"`

	UserId string `json:"userId"`
}

type Session struct {
	SessionId string `json:"sessionId"`
	CourseId  string `json:"courseId"`

	Completed bool `json:"completed"`

	SessionScore float64 `json:"sessionScore"` // set to 1 for 100% and max XP lmao
	AverageScore float64 `json:"averageScore"`

	StartingCourseProficiency float64 `json:"startingCourseProficiency"`
	StartingProficiency       int     `json:"startingProficiency"`

	EndingCourseProficiency float64 `json:"endingCourseProficiency"`
	EndingCourseScore       float64 `json:"endingCourseScore"`
	EndingProficiency       float64 `json:"endingProficiency"`

	ModulesCorrect   int `json:"modulesCorrect"`
	ModulesGaveUp    int `json:"modulesGaveUp"`
	ModulesIncorrect int `json:"modulesIncorrect"`
	ModulesStudied   int `json:"modulesStudied"`
	ModulesTested    int `json:"modulesTested"`

	Options struct {
		HasHardestQuestionContent bool `json:"hasHardestQuestionContent"`
	} `json:"options"`

	SectionIds []string `json:"sectionIds"`
	ContentIds []string `json:"contentIds"`

	SessionType string `json:"sessionType"`

	TimeStarted  time.Time `json:"timeStarted"`
	TimeFinished time.Time `json:"timeFinished"`
}
