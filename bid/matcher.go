package bid 

import (
	"fmt"
	"net/http"
	"io"
	"time"
	"math/rand"
	"github.com/valyala/fastjson"
	"encoding/json"
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sahilsp22/mini-bidder/metrics"
	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/utils"
	"github.com/sahilsp22/mini-bidder/logger"

)

// Matcher is responsible for handling Bid request and sending appropriate response
type Matcher struct{
	pg 			*db.PgClient
	mc 			*db.MCacheClient
	logger 		*logger.Logger

	// prometheus metrics and registry
	reg 		*prometheus.Registry
	bidreqs 	metrics.Counter
	bidresp 	metrics.Counter
	bidperadv 	*prometheus.CounterVec
}

var jsonpool = &fastjson.ParserPool{}

// Returns a matcher object
func NewMatcher(pg *db.PgClient, mc *db.MCacheClient, logger *logger.Logger) *Matcher {

	reg := prometheus.NewRegistry()
	// counts total bid requests recieved
	brq := metrics.NewCounter(metrics.Opts{        
		Name: "bid_requests_recieved",
		Help: "The total number of bid requests recieved",
	})
	
	// counts total bid responses sent zero and non-zero both
	brsp := metrics.NewCounter(metrics.Opts{
		Name: "bid_responses_sent",
		Help: "The total number of bid responses sent",
	})
	
	// counts total bid responses sent per advertiser
	bpadv := metrics.NewCounterVec(metrics.Opts{
		Name: "bid_response_per_advertiser",
		Help: "Total bid responses sent per advertiser",
	},[]string{"advertiserID"})
	
	reg.MustRegister(brq)
	reg.MustRegister(brsp)
	reg.MustRegister(bpadv)
	return &Matcher{
		pg:pg,
		mc:mc,
		logger:logger,
		reg:	   reg,
		bidreqs:   brq,
		bidresp:   brsp,
		bidperadv: bpadv,
	}
}

// Handler for incoming Bid request
func (m *Matcher) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request on: ",r.URL)
	m.bidreqs.Inc()
	// fmt.Println(r.Header)
	// fmt.Println("Request:\n",r)

	// Responsible for parsing bid request
	br,err := m.NewBidRequest(r)
	if err!=nil {
		m.WriteNoBidResponse(w,br,2)
		m.logger.Print(err,", Sent No Bid response")
		return
	}
	err = br.validate()
	if err != nil {
		m.WriteNoBidResponse(w,br,2)
		m.logger.Print(err,", Sent No Bid response")
		return
	}

	res := m.EvaluateBidRequest(w,br)
	if res==nil {
		return
	}
	
	// return 0 bid if no seatbids
	if res.SeatBid == nil {
		w.WriteHeader(http.StatusNoContent)
		m.logger.Print("No matching creatives, Sent No Content Response")
		return
	}
	
	// fmt.Println(res)
	b,err := json.Marshal(res)
	if err!=nil {
		m.logger.Print("Erro marshaling non-zero bid respons, Sent No Content Response,")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// update budget and metrics
	go func(*BidResponse){
		m.updateBudget(*res)
		m.updateMetrics(*res)
	}(res)

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
	
	m.logger.Print("Bid Response sent : ",res)
	return
}

// Parse the bid request and reurn a BidRequest object
func(m *Matcher) NewBidRequest(r *http.Request) (*BidRequest,error) {
	body,err := io.ReadAll(r.Body)
	if err!=nil {
		return nil,err
	}
	err = r.Body.Close()
	if err!=nil {
		m.logger.Print("Error closing request body")
	}
	// fmt.Println("Body:\n",string(body))
	p:=jsonpool.Get()
	v,err :=p.ParseBytes(body)
	if err !=nil {
		return nil,err
	}
	
	var (
		id string
		at int
		site *fastjson.Value
		device *fastjson.Value
		publisher *fastjson.Value
		geo *fastjson.Value
		user *fastjson.Value
	)
	
	id = string(v.GetStringBytes("id"))
	at = v.GetInt("at")
	site = v.Get("site")
	device = v.Get("device")
	publisher = site.Get("publisher")
	geo = device.Get("geo")
	user = device.Get("user")
	
	var br *BidRequest
	br = &BidRequest{}
	rawImps := v.GetArray("imp")
	br.Imps = make([]*Impression,len(rawImps))
	for i:=0;i<len(rawImps);i++ {
		// fmt.Println("in loop")
		raw := rawImps[i]
		var imp *Impression
		imp = &Impression{}
		var media *fastjson.Value
		var mediaType string
		// impid,err := strconv.Unquote(raw.Get("id").String())
		// if err!=nil {
			
		// }
		impid := string(raw.GetStringBytes("id"))

		imp.ID = impid
		imp.TagID = raw.Get("tagid").String()
		imp.Secure = raw.GetInt("secure")
		if raw.Exists("banner") {
			mediaType = "banner"
			media = raw.Get("banner")
		} else if raw.Exists("video") {
			mediaType = "video"
			media = raw.Get("video")
		} else if raw.Exists("audio") {
			mediaType = "audio"
			media = raw.Get("audio")
		} else if raw.Exists("native") {
			mediaType = "native"
			media = raw.Get("native")
		}

		imp.MediaType = mediaType
		imp.W = media.GetInt("w")
		imp.H = media.GetInt("h")

		br.Imps[i]=imp
	}

	br.ID = id

	if at == 0 {
		at = 2
	}
	br.At = at

	if site != nil {
		br.SiteID = string(site.GetStringBytes("id"))
		br.Domain = string(site.GetStringBytes("domain"))
	}
	
	if publisher != nil {
		br.PublisherID = string(publisher.GetStringBytes("id"))
		br.PublisherName = string(publisher.GetStringBytes("name"))
	}

	if device != nil {
		br.DeviceType = device.GetInt("id")
	}

	if geo != nil {
		br.Country = string(geo.GetStringBytes("country"))
		br.Region = string(geo.GetStringBytes("region"))
	}

	if user != nil {
		br.UserID = string(user.GetStringBytes("id"))
	}

	return br,nil
}

// Evalutes bid request to find matching bids
func (m *Matcher) EvaluateBidRequest(w http.ResponseWriter, br *BidRequest) *BidResponse {
	// if br.Imps == nil {
	// 	m.WriteNoBidResponse(w,br,2)
	// 	return nil
	// }

	creatives,err := m.GetAllCreatives()
	if err!=nil {
		m.logger.Print("Could not read creatives. Sent No Content Response")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	budgetMap := make(map[string]*utils.Budget)

	m.GetBudgets(budgetMap,creatives)

	var seatBids []SeatBid
	for _,imp := range br.Imps {
		price:=0.0
		var seat string
		for _,crtv := range creatives {

			crH,err := strconv.Atoi(crtv.Height)
			if err!=nil {
				continue
			}
			crW,err := strconv.Atoi(crtv.Width)
			if err!=nil {
				continue
			}

			// Matching impression to creatives
			if imp.MediaType == crtv.AdType && imp.W == crW && imp.H == crH {
				// skip if no advertiser found
				budget,ok := budgetMap[crtv.AdvertiserID]
				if !ok {
					continue
				}
				// skip if theres is no budget
				rem,err :=strconv.ParseFloat(budget.RemBudget,64)
				if err!=nil || rem==0 {
					continue
				}

				cpm,err :=strconv.ParseFloat(budget.CPM,64)
				if err!=nil {
					continue
				}
				// skip if budget is less
				if rem >= cpm/1000 {
					if cpm > price {
						price = cpm
						seat = budget.AdvID
					}
				} else {
					continue
				}
			} else {
				continue
			}
		}
		// add setabid only if price exists
		if price > 0 {
			rand := rand.New(rand.NewSource(time.Now().UnixNano()))
			var bid Bid
			bid.ID = strconv.Itoa(rand.Intn(1e9))
			bid.ImpID = imp.ID
			bid.Price = price
			bid.W = imp.W
			bid.H = imp.H

			seatBids = append(seatBids,SeatBid{Seat:seat,Bid:[]Bid{bid}})
		}
	}

	return &BidResponse{
		ID: br.ID,
		SeatBid: seatBids,
		NBR: NBR_DEFAULT,
	}
}

// retrieve all creatives from postgres
func (m *Matcher) GetAllCreatives() ([]*utils.Creative, error) {
	utils.CreativeLock()
	rows,err := m.pg.Query(context.Background(), utils.ALL_CREATIVES_QUERY)
	utils.CreativeUnLock()
	if err!=nil {
		return nil,fmt.Errorf("error reading creatives from Postgres : %v",err)
	}
	var creatives []*utils.Creative
	for rows.Next() {
		var crtv utils.Creative
		err = rows.Scan(&crtv.AdID, &crtv.Height, &crtv.Width, &crtv.AdType, &crtv.CreativeDetails, &crtv.AdvertiserID)
		if err != nil {
			return nil,fmt.Errorf("error scanning Creative rows: %v",err)
		}
		// fmt.Println(crtv)
		creatives = append(creatives,&crtv)
	}	
	if err = rows.Err(); err != nil {
		return nil,fmt.Errorf("error scanning Creative rows: %v",err)
	}
	defer rows.Close()
	return  creatives,nil
}

// create budget map for all AdvertiserIDs
func (m *Matcher) GetBudgets(bmap map[string]*utils.Budget, creatives []*utils.Creative) {
	for _,crtv := range creatives {
		var budget utils.Budget
		err := m.mc.Get(crtv.AdvertiserID,&budget)
		if err!=nil {
			m.logger.Print(err)
			continue
		}
		bmap[crtv.AdvertiserID] = &budget
	}
}

// returns a no bid response
func (m *Matcher) WriteNoBidResponse(w http.ResponseWriter, br *BidRequest, nbr int) {
	res := BidResponse{
		ID: br.ID,
		SeatBid: []SeatBid{},
		NBR: nbr,
	}

	b,err := json.Marshal(res)
	if err!=nil {
		m.logger.Print("Error Marshaling No Bid Response, Sent No Content Response")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.Write(b)
	m.logger.Print("Response sent with NBR: ",nbr)
	m.bidresp.Inc()
	return
}

// Returns a handler for Prometheus metrics
func (m *Matcher) Metrics() http.Handler {
	return promhttp.HandlerFor(m.reg,promhttp.HandlerOpts{})
}

// Update budget of advertisers for all bids sent
func (m *Matcher) updateBudget(res BidResponse){
	controller,_ := utils.NewController(m.pg,m.mc,m.logger)
	for _,sb := range res.SeatBid {
		err:=controller.UpdateAdvBudget(sb.Seat)
		if err!=nil {
			m.logger.Fatal(err)
		}
		m.bidperadv.WithLabelValues(sb.Seat).Inc()
		m.logger.Print("Updated budget for: ",sb.Seat)
	}
	m.bidresp.Inc()
	return
}

// update prometheus metrics
func (m *Matcher) updateMetrics(res BidResponse){
	for _,sb := range res.SeatBid {
		m.bidperadv.WithLabelValues(sb.Seat).Inc()
	}
	m.bidresp.Inc()
	return
}