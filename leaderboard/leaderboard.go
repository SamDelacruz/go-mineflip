package leaderboard

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func Post(userID string, name string, score int) {
	url := os.Getenv("LEADERBOARD_API")
	var jsonStr = []byte(fmt.Sprintf(
		`{"uid":"%s","name":"%s","score":%d}`, userID, name, score,
	))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fields := log.Fields{
			"err":    err,
			"userID": userID,
			"name":   name,
			"score":  score,
		}
		log.WithFields(fields).Error("Error posting score")
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fields := log.Fields{
		"status":  resp.Status,
		"headers": resp.Header,
		"body":    string(body),
	}
	log.WithFields(fields).Info("Score post response")
}
