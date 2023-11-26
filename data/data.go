package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	AocURL = "https://adventofcode.com/2023/leaderboard/private/view/"
)

func GetData(boardId string) (*Data, error) {
	b, err := os.ReadFile("./" + boardId + ".json")
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
	req, err := http.NewRequest("GET", AocURL+boardId+".json", nil)
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
	if err = os.WriteFile(writePath+".json", body, 0777); err != nil {
		return err
	}

	return nil
}
