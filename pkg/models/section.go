package models

type Section struct {
	Id       string `json:"id"`
	ParentId string `json:"parentId"`

	Title  string `json:"title"`
	Number string `json:"number"`

	SectionIds []string `json:"sectionIds"`
	ModuleIds  []string `json:"moduleIds"`
	ContentIds []string `json:"contentIds"`

	Contents []AnswerModule `json:"contents"` // FIXME: This should be a different struct called ContentModule as AnswerModule and ContentModule are different
}
