package models

import "time"

type AnswerModule struct {
	ModuleId  string `json:"moduleId"`
	ContentId string `json:"contentId"`
	CourseId  string `json:"courseId"`
	SectionId string `json:"sectionId"`
	SessionId string `json:"sessionId"`

	ModuleOrder int    `json:"moduleOrder"`
	ModuleType  string `json:"moduleType"`

	Completed bool `json:"completed"`
	GaveUp    bool `json:"gaveUp"`

	Content []struct{} `json:"content"`

	Score       int               `json:"score"`
	ModuleScore AnswerModuleScore `json:"moduleScore"`
	Submitted   bool              `json:"submitted"`

	TestingActive bool `json:"testingActive"`

	TimeFinished time.Time `json:"timeFinished"`
	TimeStarted  time.Time `json:"timeStarted"`
}

type AnswerModuleScore struct {
	Score int `json:"score"`
}

type ContentModule struct {
	/*
			EXAMPLE DATA
			            "courseId": "ddff7b40-4794-11e8-840f-39fdc9615de8",
	                    "id": "3ab631a0-4c7a-11e8-b481-f529820d2afa",
	                    "moduleDifficulty": 1,
	                    "moduleType": "concept",
	                    "parentId": "7e4f7350-4c79-11e8-b481-f529820d2afa"
	*/

	Id       string `json:"id"`
	ParentId string `json:"parentId"`
	CourseId string `json:"courseId"`

	ModuleDifficulty int    `json:"moduleDifficulty"`
	ModuleType       string `json:"moduleType"`
}
