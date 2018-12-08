package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rcsubra2/burrito/src/db"
	"github.com/rcsubra2/burrito/src/mockredis"
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
	client db.Database
}

// NewBurritoServer  create server server, and initialize route handlers
func NewBurritoServer(
	rts *parser.ParsedRoutes,
	mockData map[string]string,
) (*BurritoServer, error) {
	r := handler.NewRouter()
	var cli db.RedisDBInterface
	if mockData == nil {
		cli = redis.NewClient(&redis.Options{
			Addr:     "localhost:9000",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	} else {
		cli = mockredis.NewMockRedisClient(mockData)
	}
	server := &BurritoServer{
		router: r,
		client: db.NewRedisDB(cli),
	}
	for k, methodMap := range rts.Routes {
		server.addHandler(k, methodMap)
	}
	err := server.router.CheckRoutes()
	return server, err
}

func (bs *BurritoServer) queryDB(res parser.Resp, group *environment.EnvironmentGroup) {
	data, _ := res.DBReq.Run(bs.client, *group)
	for k, v := range data {
		entry := *environment.CreateStringEntry(k, v)
		group.Resp.Add(entry)
	}

}


func (bs *BurritoServer) render(
	res parser.Resp,
	w http.ResponseWriter,
	r *http.Request,
	group *environment.EnvironmentGroup,
) bool {
	if res.RespType == "FILE" {
		w.Header().Set("Content-type", "text/html")
		f, err := template.ParseFiles(res.Body.(string))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
		}
		d := group.Dump()
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
		bs.queryDB(res, group)
		return true
	}
	return false
}

// renderChain - render the list of Resp objects, until a data response is sent
func (bs *BurritoServer) renderChain(
	responses []parser.Resp,
	w http.ResponseWriter,
	r *http.Request,
	group *environment.EnvironmentGroup,
) {
	chainContinue := true
	for i, res := range responses {
		chainContinue = bs.render(res, w, r, group)
		if !chainContinue {
			if i != len(responses)-1 {
				log.Println("WARN: Response sent before all actions completed!")
			}
			break
		}
	}
	if chainContinue {
		WriteJson(w, group.Resp.Data())
	}
}



// addHandler - for given path and method map
func (bs *BurritoServer) addHandler(path string, methodMap map[string][]parser.Resp) {
	bs.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request, env * environment.Env) {
		// set form data and url data in an environment
		r.ParseForm()
		formEnv := environment.CreateEnv()
		for k, _ := range r.Form {
			entry := environment.CreateStringEntry(k, r.Form.Get(k))
			formEnv.Add(*entry)
		}
		// set response data in an environment
		respEnv := environment.CreateEnv()


		group := environment.CreateEnvironmentGroup(env, formEnv, respEnv)

		for k, v := range methodMap {
			if r.Method == k {
				bs.renderChain(v, w, r, group)
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