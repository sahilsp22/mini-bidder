package server

import (
	"fmt"
	// "net"
	"net/http"

	"github.com/sahilsp22/mini-bidder/logger"
)

var srvlog *logger.Logger

func init() {
	srvlog = logger.InitLogger(logger.SERVER)
}

type Server struct{
	srvr http.Server
	handler http.Handler
}

type Route struct{
	Path string
	Handler http.HandlerFunc
}

func (s *Server) AddRoutes(routes []Route) {
	mux := http.NewServeMux()

	for _,rt := range routes {
		mux.HandleFunc(rt.Path,rt.Handler)
		srvlog.Printf("Adding route: %s",rt.Path)
	}
	s.handler = mux
}

func (s *Server) Listen(port int){
	p:=fmt.Sprintf(":%v",port)
	s.srvr = http.Server{
		Addr: p,
		Handler: s.handler,
	}

	err:= s.srvr.ListenAndServe()
	if err!=nil || err == http.ErrServerClosed{
		srvlog.Fatal("Error starting server :",err)
	}
	srvlog.Printf("Server running on port :%s",port)
}
