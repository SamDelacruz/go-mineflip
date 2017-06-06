package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmcvetta/randutil"
	"github.com/samdelacruz/go-mineflip/leaderboard"
	"github.com/samdelacruz/go-mineflip/model"
)

var games = make(map[string]*model.Game)

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := randutil.AlphaString(5)
	g := model.Game{ID: id, Board: model.GenBoard()}
	games[id] = &g
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
		// Early return if game has finished
		if g.Won() || g.Lost() {
			w.WriteHeader(http.StatusOK)
			data, _ := g.ToJSON()
			w.Write(data)
			return
		}
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
		g.AddMove(i)
		if g.Won() {
			// Post the score to player's leaderboard
			go func() {
				token := parseJWT(r)
				if token != nil {
					userID := token["sub"].(string)
					givenName := token["given_name"].(string)
					if len(userID) > 0 && len(givenName) > 0 {
						leaderboard.Post(userID, givenName, g.GetScore())
					}
				}
			}()
		}
		w.WriteHeader(http.StatusOK)
		data, _ := g.ToJSON()
		w.Write(data)
		return
	}
	http.NotFound(w, r)
}

func getGame(r *http.Request) (*model.Game, bool) {
	vars := mux.Vars(r)
	id := vars["id"]
	g, ok := games[id]
	return g, ok
}
