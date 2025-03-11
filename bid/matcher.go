package bid 

import (
	"fmt"
	"net/http"
	"io"

	// "github.com/sahilsp22/mini-bidder/logger"
	"github.com/sahilsp22/mini-bidder/metrics"
)

type Matcher struct{
	bidmetric metrics.Counter
}

func NewMatcher() *Matcher {
	metric := metrics.NewCounter(metric.Opts{
		Name: "bid_requests_recieved",
		Help: "The total number of bid requests recieved",
	})

	return &Matcher{bidmetric:metric}
}

func (m *Matcher) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved request on: ",r.URL)
	fmt.Println("Request:\n",r)
	body,_ := io.ReadAll(r.Body)
	fmt.Println("Body:\n",string(body))

	w.Write([]byte("Bid Request Recieved!!"))
	m.bidmetric.Inc()
	// io.WriteString(w,"Bid Request Recieved!!")
	return
}