package bid

import(
	// "encoding/json"
)

type BidRequest struct {
	Id string `json:"id"`
	Imps *[]Impression `json:"imps"`
	At int `json:"at,omitempty"`
	Site *Site `json:"site,omitempty"`
	Device *Device `json:"device",omitempty`
}

type Impression struct{
	Id string `json:"id"`
	TagId string `json"tagid",omitempty`
	Secure string `json:"secure",omitempty`
	Banner *Banner `json:"banner",omitempty`
}

type Banner struct{
	W int32 `json:"w"`
	H int32 `json:"h"`
}

type Site struct{
	Id string `json:"id",omitempty`
	Publisher *Publisher `json:"publisher",omitempty`
	Domain string `json:"domain",omitempty`
}

type Publisher struct{
	Id string `json:"id",omitempty`
	Name string `json:"name",omitempty`
}

type Device struct{
	Geo *Geo `json:"geo",omitempty`
	DeviceType int `json:"devicetype",omitempty`
	User *User `json:"user",omitempty`
}

type Geo struct{
	Country string `json:"country",omitempty`
	Region string `json:"region",omitempty`
}

type User struct{
	Id string `json:"id",omitempty`
}