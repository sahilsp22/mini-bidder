package bid

const (
	NBR_DEFAULT = -1		// default nbr code
)

type BidResponse struct {
	ID 		string			`json:"id"`
	SeatBid []SeatBid		`json:"seatbid"`
	NBR 	int				`json:"nbr",omitempty`
}

type SeatBid struct {
	Seat 	string			`json:"seat"`
	Bid 	[]Bid			`json:"bid"`
}

type Bid struct {
	ID 		string			`json:"id"`
	ImpID 	string			`json:"impid"`
	Price 	float64			`json:"price"`
	W 		int				`json:"w",omitempty`
	H 		int				`json:"h",omitempty`
}
