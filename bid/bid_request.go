package bid

import(
	// "encoding/json"
	"fmt"
)

type BidRequest struct {
	ID string 
	Imps []*Impression 
	At int 

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
	ID string 
	TagID string 
	Secure int 
	MediaType string 
	W int
	H int
}

type Media struct{
	W int `json:"w"`
	H int `json:"h"`
}

// Validates a bid request
func (br *BidRequest) validate() error {
	if br.ID == "" {
		return fmt.Errorf("Bid request ID is nil")
	}
	
	if len(br.Imps)==0 {
		return fmt.Errorf("no Impressions in Bid request")
	}
	
	// check if any impression has missing id
	for _,imps := range br.Imps {
		if imps.ID == ""{
			return fmt.Errorf("impression ID is nil")
		}
	}
	
	return nil
}




type Site struct{
	Id string 
	Publisher *Publisher 
	Domain string 
}

type Publisher struct{
	Id string 
	Name string 
}

type Device struct{
	Geo *Geo 
	DeviceType int 
	User *User 
}

type Geo struct{
	Country string 
	Region string 
}

type User struct{
	Id string 
}