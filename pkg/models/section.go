package models

type Section struct {
	Id       string `json:"id"`
	ParentId string `json:"parentId"`

	Title  string `json:"title"`
	Number string `json:"number"`

	SectionIds []string `json:"sectionIds"`
	ModuleIds  []string `json:"moduleIds"`
	ContentIds []string `json:"contentIds"`

	Contents []SectionContent `json:"contents"`
}

type SectionContent struct {
	Id             string          `json:"id"`
	ContentModules []ContentModule `json:"contentModules"`
}
