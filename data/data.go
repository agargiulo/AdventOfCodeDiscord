package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	aocURLBase = "https://adventofcode.com/%d/"
)

var (
	dataDir    = os.Getenv("DATA_DIR")
	curYear    = time.Now().Year()
	aocBaseUrl = fmt.Sprintf(aocURLBase, curYear)
)

var (
	LearderBoardTitle = fmt.Sprintf("🎄 %d Leaderboard 🎄", curYear)
)

func GetData(boardId string) (*Data, error) {
	b, err := os.ReadFile(dataDir + boardId + ".json")
	if err != nil {
		return nil, err
	}

	var D Data
	if err = json.Unmarshal(b, &D); err != nil {
		return nil, err
	}
	return &D, nil
}

func FetchData(boardId, sessionToken, writePath string) error {
	// Form request to adventofcode API
	req, err := http.NewRequest("GET", AocLeaderboardUrl(boardId)+".json", nil)
	if err != nil {
		return err
	}

	// Add session token
	req.Header.Add("Cookie", "session="+sessionToken)

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("error fetching data from leaderboard: %s", boardId)
	}

	// Write data to file
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = os.WriteFile(dataDir+writePath+".json", body, 0777); err != nil {
		return err
	}

	return nil
}

func AocLeaderboardUrl(boardId string) string {
	aocLeaderboardURL, err := url.JoinPath(aocBaseUrl, "leaderboard/private/view/", boardId)
	if err != nil {
		panic(err)
	}
	return aocLeaderboardURL
}

func ProblemUrl(day int) string {
	problemUrl, err := url.JoinPath(aocBaseUrl, "day", strconv.Itoa(day))
	if err != nil {
		panic(err)
	}
	return problemUrl
}
