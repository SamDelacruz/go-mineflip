package api

import (
	"github.com/gorilla/mux"
	"github.com/jmcvetta/randutil"
	"github.com/samdelacruz/go-mineflip/model"
	"net/http"
	"strconv"
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
	if g, ok := getGame(r); ok {
		data, _ := g.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
		return
	}
	http.NotFound(w, r)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	if g, ok := getGame(r); ok {
		vars := mux.Vars(r)

		x, err := strconv.Atoi(vars["x"])
		if err != nil || x > 4 || x < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		y, err := strconv.Atoi(vars["y"])

		if err != nil || y > 4 || y < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i := y*5 + x
		if moveAvail(i, g) {
			newGame := model.Game{ID: g.ID, Moves: append(g.Moves, i), Board: g.Board}
			games[g.ID] = newGame
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.NotFound(w, r)
}

func moveAvail(i int, g model.Game) bool {
	for _, m := range g.Moves {
		if m == i {
			return false
		}
	}
	return true
}

func getGame(r *http.Request) (model.Game, bool) {
	vars := mux.Vars(r)
	id := vars["id"]
	g, ok := games[id]
	return g, ok
}
