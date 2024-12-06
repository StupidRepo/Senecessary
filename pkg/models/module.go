package models

type Module struct {
	ModuleId  string `json:"moduleId"`
	ContentId string `json:"contentId"`
	CourseId  string `json:"courseId"`
	SectionId string `json:"sectionId"`
	SessionId string `json:"sessionId"`

	Completed bool `json:"completed"`
}
