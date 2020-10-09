package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type point struct {
	X, Y int
}

// a server shows the current state of the TSP via an HTTP server.
type server struct {
	solutions chan []point // channel type is temporary (should probably be a channel of Individual)
}

func newServer() *server {
	return &server{
		solutions: make(chan []point),
	}
}

func (s *server) serve(host string) {
	log.Printf("Server starting, point your browser to http://%s\n", host)

	http.Handle("/", http.FileServer(http.Dir("./example/tsp")))
	http.HandleFunc("/ws", s.ws)

	go http.ListenAndServe(host, nil)
}

func (s *server) start() error {
	tick := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			s.solutions <- randomPath(4)
		}
	}
}

const xmax, ymax = 200, 200

var updates int

func randomPath(length int) []point {
	updates++
	if updates%2 == 0 {
		path := make([]point, length)
		path[0] = point{X: 0, Y: 0}
		path[1] = point{X: 0, Y: ymax}
		path[2] = point{X: xmax, Y: ymax}
		path[3] = point{X: xmax, Y: 0}
		return path
	}

	path := make([]point, length)
	for i := range path {
		path[i].X = rand.Intn(xmax)
		path[i].Y = rand.Intn(ymax)
	}

	return path
}

func (s *server) ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	fmt.Println("ws start")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	if err = s.runTSP(ws); err != nil {
		log.Println(err)
	}

	fmt.Println("ws end")
}

type initMessage struct {
	Width, Height int
}

func (s *server) runTSP(conn *websocket.Conn) error {
	conn.WriteJSON(initMessage{Width: xmax, Height: ymax})

	for sol := range s.solutions {
		if err := conn.WriteJSON(sol); err != nil {
			return err
		}
	}

	return nil
}

// Perm calls f with each permutation of a.
func Perm(a []int, f func([]int)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, f func([]int), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

/*

	"github.com/arl/statsviz/websocket"
)

func init() {
	http.Handle("/debug/statsviz/", Index)
	http.HandleFunc("/debug/statsviz/ws", Ws)
}

// Index responds to a request for /debug/statsviz with the statsviz HTML page
// which shows a live visualization of the statistics sent by the application
// over the websocket handler Ws.
//
// The package initialization registers it as /debug/statsviz/.
var Index = http.StripPrefix("/debug/statsviz/", http.FileServer(assets))

// Ws upgrades the HTTP server connection to the WebSocket protocol and sends
// application statistics every second.
//
// If the upgrade fails, an HTTP error response is sent to the client.
// The package initialization registers it as /debug/statsviz/ws.
func Ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("can't upgrade HTTP connection to Websocket protocol:", err)
		return
	}
	defer ws.Close()

	err = sendStats(ws)
	if err != nil {
		log.Println(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}

const defaultSendPeriod = time.Second

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn) error {
	tick := time.NewTicker(defaultSendPeriod)

	var stats stats
	for {
		select {
		case <-tick.C:
			runtime.ReadMemStats(&stats.Mem)
			stats.NumGoroutine = runtime.NumGoroutine()
			if err := conn.WriteJSON(stats); err != nil {
				return err
			}
		}
	}
}
*/
