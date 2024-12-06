package models

type Provider struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`

	ProviderId string `json:"providerId"`
	Uid        string `json:"uid"`
}
