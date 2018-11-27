package server

import (
	"encoding/json"
	"fmt"
	"github.com/rcsubra2/burrito/src/db"
	"html/template"
	"log"
	"net/http"

	"github.com/rcsubra2/burrito/src/environment"
	"github.com/rcsubra2/burrito/src/handler"
	"github.com/rcsubra2/burrito/src/parser"
)

// BurritoServer - the webserver that serves parsed routes
type BurritoServer struct {
	router *handler.Router
	client *db.RedisClient
}

// NewBurritoServer  create server server, and initialize route handlers
func NewBurritoServer(rts *parser.ParsedRoutes) *BurritoServer {
	r := handler.NewRouter()
	server := &BurritoServer{
		router: r,
		client: db.NewRedisClient("localhost:9000"),
	}
	for k, methodMap := range rts.Routes {
		server.addHandler(k, methodMap)
	}
	return server
}

func (bs *BurritoServer) queryDB(res parser.Resp, urlEnv *environment.Env, respEnv *environment.Env) {
	req, ok := res.Body.(db.Req)
	if ok {
		if req.Method == "GET" {
			data := bs.client.Get(req.GetReq, []*environment.Env{urlEnv, respEnv})
			for k, v := range data {
				entry := *environment.CreateStringEntry(k, v)
				respEnv.Add(entry)
			}
		}
	}
}


func (bs *BurritoServer) render(
	res parser.Resp,
	w http.ResponseWriter,
	r *http.Request,
	urlEnv * environment.Env,
	respEnv * environment.Env,
) bool {
	if res.RespType == "FILE" {
		w.Header().Set("Content-type", "text/html")
		f, err := template.ParseFiles(res.Body.(string))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
		}
		d := environment.CreateBurritoTemplateData(urlEnv, respEnv)
		f.Execute(w, &d)
		return false

	} else if res.RespType == "STR" {
		w.Header().Set("Content-type", "text/html")
		fmt.Fprintf(w, res.Body.(string))
		return false
	} else if res.RespType == "JSON" {
		WriteJson(w, res.Body)
		return false
	} else if res.RespType == "DB" {
		bs.queryDB(res, urlEnv, respEnv)
		return true
	}
	return false
}

// renderChain - render the list of Resp objects, until a data response is sent
func (bs *BurritoServer) renderChain(
	responses []parser.Resp,
	w http.ResponseWriter,
	r *http.Request,
	urlEnv * environment.Env,
	respEnv * environment.Env,
) {
	chainContinue := true
	for i, res := range responses {
		chainContinue = bs.render(res, w, r, urlEnv, respEnv)
		if !chainContinue {
			if i != len(responses)-1 {
				log.Println("WARN: Response sent before all actions completed!")
			}
			break
		}
		fmt.Println(respEnv)
	}
	fmt.Println(respEnv)
	// TODO: check if no response was sent back, if so send back the respEnv as JSON
	if chainContinue {
		WriteJson(w, respEnv.Data())
	}
}



// addHandler - for given path and method map
func (bs *BurritoServer) addHandler(path string, methodMap map[string][]parser.Resp) {
	bs.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request, env * environment.Env) {
		// TODO: Create Burrito Env with all url data included, and pass that into renderChain
		respEnv := environment.CreateEnv()
		for k, v := range methodMap {
			if r.Method == k {
				bs.renderChain(v, w, r, env, respEnv)
			}
		}
	})
}

// Run - run the burrito server
func (bs *BurritoServer) Run() {
	log.Fatal(http.ListenAndServe(":8080", bs.router))
}


// Write interface as JSON to http.ResponseWriter
func WriteJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-type", "app/json")
	json, err := json.Marshal(data)
	if err != nil {
		panic("Error: json data not parsed correctly")
	} else {
		w.Write(json)
	}
}