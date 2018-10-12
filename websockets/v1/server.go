package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

// Player stores a name with a number of wins
type Player struct {
	Name string
	Wins int
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	p.playGame(conn)
}

func (p *PlayerServer) playGame(conn *websocket.Conn) {
	_, numberOfPlayersMsg, err := conn.ReadMessage()
	if err != nil {
		return
	}

	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	game := NewTexasHoldem(&WebSocketBlindAlerter{conn}, p.store)
	game.Start(numberOfPlayers)

	_, winner, err := conn.ReadMessage()
	if err != nil {
		return
	}

	game.Finish(string(winner))
}

type WebSocketBlindAlerter struct {
	*websocket.Conn
}

func (w *WebSocketBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		w.WriteMessage(1, []byte(fmt.Sprintf("The blind is now %d", amount)))
	})
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:

	}
	t, err := template.ParseFiles("../game.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
	w.Header().Add("content-type", "text/html")
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
