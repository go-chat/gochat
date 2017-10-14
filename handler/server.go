package handler

import (
	"log"
	"net/http"

	"github.com/go-chat/gochat/config"
	"github.com/go-chat/gochat/store"
	"github.com/gorilla/mux"
	websocket "github.com/gorilla/websocket"

	_ "github.com/lib/pq"
)

type Server struct {
	Store    *store.Store
	Router   *mux.Router
	Clients  map[*http.Request]*websocket.Conn
	Upgrader websocket.Upgrader
}

// global server
var Srv *Server

func NewServer(cfg *config.Config) {
	Srv = &Server{}

	Srv.Store = store.NewStore(cfg)

	router := NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../static/"))))

	Srv.Router = router

	Srv.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Srv.Clients = make(map[*http.Request]*websocket.Conn)
}

func Serve(cfg *config.Config) {
	log.Printf("Listen and serve at :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, Srv.Router))
}

func StopServer() {
	Srv.Store.Close()
}
