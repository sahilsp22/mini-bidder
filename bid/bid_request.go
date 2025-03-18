package bid

import(
	// "encoding/json"
	"fmt"
)

type BidRequest struct {
	ID string `json:"id"`
	Imps []*Impression `json:"imps"`
	At int `json:"at,omitempty"`

	// Site Detials
	SiteID string
	Domain string
	PublisherID string
	PublisherName string

	// Device Details
	DeviceType int
	Country string
	Region string
	UserID string
}

type Impression struct{
	ID string `json:"id"`
	TagID string `json"tagid",omitempty`
	Secure int `json:"secure",omitempty`
	MediaType string 
	W int
	H int
}

type Media struct{
	W int `json:"w"`
	H int `json:"h"`
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

func (br *BidRequest) validate() error {
	if br.ID == "" {
		return fmt.Errorf("Bid request ID is nil")
	}

	if len(br.Imps)==0 {
		return fmt.Errorf("No Impressions in Bid request")
	}

	return nil
}