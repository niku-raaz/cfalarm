package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CFProblem struct {
	ContestID int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
	Rating    int    `json:"rating"`
}

func FetchCFProblemsByRating(rating int) ([]CFProblem, error) {
	url := fmt.Sprintf("https://codeforces.com/api/problemset.problems?tags=%d", rating)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Result struct {
			Problems []CFProblem `json:"problems"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Result.Problems, nil
}
