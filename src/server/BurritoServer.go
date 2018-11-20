package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rcsubra2/burrito/src/mux"
	"github.com/rcsubra2/burrito/src/parser"
)

// BurritoServer - the webserver that serves parsed routes
type BurritoServer struct {
	router *mux.Router
}

// NewBurritoServer  create server server, and initialize route handlers
func NewBurritoServer(rts *parser.ParsedRoutes) *BurritoServer {
	r := mux.NewRouter()
	server := &BurritoServer{
		router: r,
	}
	for k, methodMap := range rts.Routes {
		server.addHandler(k, methodMap)
	}
	return server
}

func (bs *BurritoServer) render(res parser.Resp, w http.ResponseWriter, r *http.Request) bool {
	if res.RespType == "FILE" {
		w.Header().Set("Content-type", "text/html")
		f, err := ioutil.ReadFile(string(res.Body))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
		}
		w.Write(f)
		return false
	} else if res.RespType == "STR" {
		// TODO
		w.Header().Set("Content-type", "text/html")
		fmt.Fprintf(w, res.Body)
		return false
	} else if res.RespType == "CONT" {
		// TODO
		return true
	}
	return false
}

// renderChain - render the list of Resp objects, until a data response is sent
func (bs *BurritoServer) renderChain(responses []parser.Resp, w http.ResponseWriter, r *http.Request) {
	for i, res := range responses {
		isRedirect := bs.render(res, w, r)
		if !isRedirect {
			if i != len(responses)-1 {
				log.Println("WARN: Response sent before all actions completed!")
			}
			break
		}
	}
}

// addHandler - for given path and method map
func (bs *BurritoServer) addHandler(path string, methodMap map[string][]parser.Resp) {
	bs.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		for k, v := range methodMap {
			if r.Method == k {
				bs.renderChain(v, w, r)
			}
		}
	})
}

// Run - run the burrito server
func (bs *BurritoServer) Run() {
	log.Fatal(http.ListenAndServe(":8080", bs.router))
}
