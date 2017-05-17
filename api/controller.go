package api

import (
	"github.com/gorilla/mux"
	"github.com/jmcvetta/randutil"
	"github.com/samdelacruz/go-mineflip/model"
	"net/http"
)

var games = make(map[string]model.Game)

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := randutil.AlphaString(5)
	g := model.Game{ID: id, Board: model.GenBoard()}
	games[id] = g
	data, _ := g.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
}

func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if g, ok := games[id]; ok {
		data, _ := g.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
		return
	}
	http.NotFound(w, r)
}
