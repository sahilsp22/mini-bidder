package bid 

import (
	"fmt"
	"net/http"
	"io"
	"time"
	"math/rand"
	"github.com/valyala/fastjson"

	// "github.com/sahilsp22/mini-bidder/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sahilsp22/mini-bidder/metrics"

)

type Matcher struct{
	reg 		*prometheus.Registry
	bidreqs 	metrics.Counter
	bidresp 	metrics.Counter
	bidperadv 	*prometheus.CounterVec
}

var jsonpool = &fastjson.ParserPool{}

func NewMatcher() *Matcher {

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
		reg:	   reg,
		bidreqs:   brq,
		bidresp:   brsp,
		bidperadv: bpadv,
	}
}

func (m *Matcher) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request on: ",r.URL)
	// fmt.Println(r.Header)
	// fmt.Println("Request:\n",r)
	_ = m.NewBidRequest(r)
	
	w.Write([]byte("Bid Request Recieved!!"))
	m.bidreqs.Inc()
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rand.Intn(10) > 2 {
		m.bidresp.Inc()
		ad := fmt.Sprintf("advtest%v",rand.Intn(4)+1)
		m.bidperadv.WithLabelValues(ad).Inc()
		w.Write([]byte(fmt.Sprintf("Bid Response Sent for advertiser : %v",ad)))
		return
	}
	w.Write([]byte("No Bid Response"))

	return
}

func(m *Matcher) NewBidRequest(r *http.Request) *BidRequest{
	body,_ := io.ReadAll(r.Body)
	fmt.Println("Body:\n",string(body))
	p:=jsonpool.Get()
	v,_:=p.ParseBytes(body)
	fmt.Println(v.Get("site"))
	return &BidRequest{}
}

func (m *Matcher) Metrics() http.Handler {
	return promhttp.HandlerFor(m.reg,promhttp.HandlerOpts{})
}