package models

import "time"

type User struct {
	AvatarUrl string `json:"avatarUrl"`
	Birthdate string `json:"birthdate"`

	GivenName   string `json:"givenName"`
	FamilyName  string `json:"familyName"`
	DisplayName string `json:"displayName"`

	Email        string        `json:"email"`
	EmailHistory []interface{} `json:"emailHistory"`

	EnglishAsAdditionalLanguage bool `json:"englishAsAdditionalLanguage"`

	ExternalSchoolId string `json:"externalSchoolId"`
	ExternalUserId   string `json:"externalUserId"`

	LastVisitTime time.Time `json:"lastVisitTime"`

	ManagedBy    string     `json:"managedBy"`
	ProviderData []Provider `json:"providerData"`

	Ethnicity string `json:"ethnicity"`
	Gender    string `json:"gender"`

	Gifted                  bool   `json:"gifted"`
	InLeaCare               bool   `json:"inLeaCare"`
	FreeSchoolMeals         bool   `json:"freeSchoolMeals"`
	PupilPremium            bool   `json:"pupilPremium"`
	SpecialEducationalNeeds string `json:"specialEducationalNeeds"`

	TimeCreated    time.Time `json:"timeCreated"`
	TimeLastSynced time.Time `json:"timeLastSynced"`
	TimeLinked     time.Time `json:"timeLinked"`

	Type   string `json:"type"`
	Upn    string `json:"upn"`
	UserId string `json:"userId"`

	Assignments []Assignment // Not from API, used for frontend
}
