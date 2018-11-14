package burrito

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// BurritoServer - the webserver that serves parsed routes
type BurritoServer struct {
	Routes *ParsedRoutes
}

// NewBurritoServer  create burrito server, and initialize route handlers
func NewBurritoServer(rts *ParsedRoutes) *BurritoServer {
	server := &BurritoServer{
		Routes: rts,
	}
	fmt.Println(rts)
	for k, methodMap := range server.Routes.routes {
		server.addHandler(k, methodMap)
	}
	return server
}

func (bs *BurritoServer) render(res Resp, w http.ResponseWriter, r *http.Request) bool {
	if res.respType == "FILE" {
		w.Header().Set("Content-type", "text/html")
		f, err := ioutil.ReadFile(string(res.body))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
		}
		w.Write(f)
		return false
	} else if res.respType == "STR" {
		// TODO
		w.Header().Set("Content-type", "text/html")
		fmt.Fprintf(w, res.body)
		return false
	} else if res.respType == "CONT" {
		// TODO
		return true
	}
	return false
}

// renderChain - render the list of Resp objects, until a data response is sent
func (bs *BurritoServer) renderChain(responses []Resp, w http.ResponseWriter, r *http.Request) {
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
func (bs *BurritoServer) addHandler(path string, methodMap map[string][]Resp) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		for k, v := range methodMap {
			if r.Method == k {
				bs.renderChain(v, w, r)
			}
		}
	})
}

// Run - run the burrito server
func (bs *BurritoServer) Run() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
