package leaderboard

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
