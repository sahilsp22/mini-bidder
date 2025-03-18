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

	// "github.com/sahilsp22/mini-bidder/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sahilsp22/mini-bidder/metrics"
	"github.com/sahilsp22/mini-bidder/db"
	"github.com/sahilsp22/mini-bidder/utils"
	"github.com/sahilsp22/mini-bidder/logger"

)

type Matcher struct{
	pg *db.PgClient
	mc *db.MCacheClient
	cntrl *utils.Controller
	logger *logger.Logger
	reg 		*prometheus.Registry
	bidreqs 	metrics.Counter
	bidresp 	metrics.Counter
	bidperadv 	*prometheus.CounterVec
}

var jsonpool = &fastjson.ParserPool{}

func NewMatcher(pg *db.PgClient, mc *db.MCacheClient, c *utils.Controller, logger *logger.Logger) *Matcher {

	reg := prometheus.NewRegistry()
	brq := metrics.NewCounter(metrics.Opts{
		Name: "bid_requests_recieved",
		Help: "The total number of bid requests recieved",
	})

	brsp := metrics.NewCounter(metrics.Opts{
		Name: "bid_responses_sent",
		Help: "The total number of bid responses sent",
	})

	bpadv := metrics.NewCounterVec(metrics.Opts{
		Name: "bid_response_per_advertiser",
		Help: "Total bid responses sent per advertiser",
	},[]string{"advid"})
	
	reg.MustRegister(brq)
	reg.MustRegister(brsp)
	reg.MustRegister(bpadv)
	return &Matcher{
		pg:pg,
		mc:mc,
		cntrl:c,
		logger:logger,
		reg:	   reg,
		bidreqs:   brq,
		bidresp:   brsp,
		bidperadv: bpadv,
	}
}

func (m *Matcher) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request on: ",r.URL)
	m.bidreqs.Inc()
	// fmt.Println(r.Header)
	// fmt.Println("Request:\n",r)
	br := m.NewBidRequest(r)
	err := br.validate()
	if err!=nil {
		m.WriteNoBidResponse(w,br,2)
		m.logger.Print(err,", Sent No Bid response")
		return
	}
	res := m.EvaluateBidRequest(w,br)
	if res==nil {
		return
	}
	
	if res.SeatBid == nil {
		w.WriteHeader(http.StatusNoContent)
		m.logger.Print("No matching creatives, Sent No Content Response")
		return
	}
	
	fmt.Println(res)
	b,err := json.Marshal(res)
	if err!=nil {
		m.logger.Print("Erro marshaling non-zero bid respons, Sent No Content Response,")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.Write(b)

	go func(*BidResponse){
		m.updateMetrics(*res)
	}(res)
	// w.Write([]byte("Bid Request Recieved!!"))
	// m.bidreqs.Inc()
	// rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	// if rand.Intn(10) > 2 {
	// 	m.bidresp.Inc()
	// 	ad := fmt.Sprintf("advtest%v",rand.Intn(4)+1)
	// 	m.bidperadv.WithLabelValues(ad).Inc()
	// 	w.Write([]byte(fmt.Sprintf("Bid Response Sent for advertiser : %v",ad)))
	// 	return
	// }
	// w.Write([]byte("No Bid Response"))

	return
}

func(m *Matcher) NewBidRequest(r *http.Request) *BidRequest{
	body,_ := io.ReadAll(r.Body)
	// fmt.Println("Body:\n",string(body))
	p:=jsonpool.Get()
	v,_:=p.ParseBytes(body)
	
	var (
		id string
		at int
		site *fastjson.Value
		device *fastjson.Value
		publisher *fastjson.Value
		geo *fastjson.Value
		user *fastjson.Value
	)
	
	id = v.Get("id").String()
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
		impid,err := strconv.Unquote(raw.Get("id").String())
		if err!=nil {

		}
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

	brid,err := strconv.Unquote(id)
	if err!=nil {

	}
	br.ID = brid
	br.At = at

	br.SiteID = site.Get("id").String()
	br.Domain = site.Get("domain").String()
	br.PublisherID = publisher.Get("id").String()
	br.PublisherName = publisher.Get("name").String()

	br.DeviceType = device.GetInt("id")
	br.Country = geo.Get("country").String()
	br.Region = geo.Get("region").String()
	br.UserID = user.Get("id").String()

	return br
}

func (m *Matcher) EvaluateBidRequest(w http.ResponseWriter, br *BidRequest) *BidResponse {
	if br.Imps == nil {
		m.WriteNoBidResponse(w,br,2)
		return nil
	}

	var creatives []utils.Creative
	creatives,err := m.GetAllCreatives()
	if err!=nil {
		m.logger.Print("Could not read creatives. Sent No Content Response")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	// impCrtvMap := make(map[string])
	var seatBids []SeatBid
	for _,imp := range br.Imps {
		price:=0.0
		var seat string
		for _,crtv := range creatives {
			crH,err := strconv.Atoi(crtv.Height)
			crW,err := strconv.Atoi(crtv.Width)
			if err!=nil {

			}
			if imp.MediaType == crtv.AdType && imp.W == crW && imp.H == crH {
				var budget utils.Budget
				err := m.mc.Get(crtv.AdvertiserID,&budget)
				if err!=nil {
					m.logger.Print(err)
					continue
				}
				cpm,err :=strconv.ParseFloat(budget.CPM,64)
				rem,err :=strconv.ParseFloat(budget.RemBudget,64)
				if rem < cpm {
					continue
				}
				if cpm > price {
					price = cpm
					seat = budget.AdvID
				}
			}
		}
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
		NBR: -1,
	}
}

func (m *Matcher) GetAllCreatives() ([]utils.Creative, error) {
	rows,err := m.pg.Query(context.Background(), "SELECT * FROM Creative_Details;")
	if err!=nil {
		return nil,fmt.Errorf("Error reading creatives from Postgres : %v",err)
	}
	var creatives []utils.Creative
	for rows.Next() {
		var crtv utils.Creative
		err = rows.Scan(&crtv.AdID, &crtv.Height, &crtv.Width, &crtv.AdType, &crtv.CreativeDetails, &crtv.AdvertiserID)
		if err != nil {
			return nil,fmt.Errorf("Error scanning Creative rows: %v",err)
		}
		// fmt.Println(crtv)
		creatives = append(creatives,crtv)
	}	
	if err = rows.Err(); err != nil {
		return nil,fmt.Errorf("Error scanning Creative rows: %v",err)
	}
	defer rows.Close()

	return creatives, nil
}

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
	return
}

func (m *Matcher) Metrics() http.Handler {
	return promhttp.HandlerFor(m.reg,promhttp.HandlerOpts{})
}

func (m *Matcher) updateMetrics(res BidResponse){
	m.bidresp.Inc()
	for _,sb := range res.SeatBid {
		err:=m.cntrl.UpdateAdvBudget(sb.Seat)
		if err!=nil {
			m.logger.Print(err)
		}
		m.bidperadv.WithLabelValues(sb.Seat).Inc()
	}
	return
}