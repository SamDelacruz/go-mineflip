package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/randutil"
	"github.com/samdelacruz/go-mineflip/leaderboard"
	"github.com/samdelacruz/go-mineflip/model"
)

var games = make(map[string]*model.Game)

func CreateGameHandler(c *gin.Context) {
	id, _ := randutil.AlphaString(5)
	g := model.Game{ID: id, Board: model.GenBoardForLevel(1)}
	games[id] = &g
	c.JSON(http.StatusCreated, g.ToRepr())
}

func GetGameHandler(c *gin.Context) {
	if g, ok := getGame(c); ok {
		c.JSON(http.StatusOK, g.ToRepr())
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func MoveHandler(c *gin.Context) {
	if g, ok := getGame(c); ok {
		// Early return if game has finished
		if g.Won() || g.Lost() {
			c.JSON(http.StatusOK, g.ToRepr())
			return
		}
		xStr := c.Param("x")
		yStr := c.Param("y")

		x, err := strconv.Atoi(xStr)
		if err != nil || x > 4 || x < 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		y, err := strconv.Atoi(yStr)

		if err != nil || y > 4 || y < 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		i := y*5 + x
		g.AddMove(i)
		if g.Won() {
			// Post the score to player's leaderboard
			go func() {
				u, exists := c.Get("user")
				if exists {
					user := u.(User)
					leaderboard.Post(user.ID, user.Name, g.GetScore())
				}
			}()
		}
		c.JSON(http.StatusOK, g.ToRepr())
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func getGame(c *gin.Context) (*model.Game, bool) {
	id := c.Param("id")
	g, ok := games[id]
	return g, ok
}
