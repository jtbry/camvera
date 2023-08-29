package web

import (
	"context"
	"log"
)

type webServer struct {
	ctx context.Context
	log *log.Logger
}

func NewWebServer(c context.Context, logger *log.Logger) *webServer {
	return &webServer{
		ctx: c,
		log: logger,
	}
}

func (w *webServer) Start() {
	// TODO
	
}
