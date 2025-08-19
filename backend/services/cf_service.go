package services

import (
	"crypto/sha512" // Import crypto package
	"encoding/hex"  // Import hex encoding
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type CFProblem struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}

// Represents a single submission from the user.status endpoint
type CFSubmission struct {
	Problem CFProblem `json:"problem"`
	Verdict string    `json:"verdict"`
}

func VerifyCFKeys(handle, key, secret string) (bool, error) {
	if handle == "" || key == "" || secret == "" {
		return false, errors.New("handle, key, and secret cannot be empty")
	}

	// 1. Choose a method and prepare parameters
	methodName := "user.info"
	params := map[string]string{
		"handles": handle,
		"apiKey":  key,
		"time":    strconv.FormatInt(time.Now().Unix(), 10),
	}

	// 2. Create the parameter string for hashing
	var paramList []string
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(paramList) // Sort params alphabetically
	paramString := ""
	for i, p := range paramList {
		if i > 0 {
			paramString += "&"
		}
		paramString += p
	}

	// 3. Create the string to be hashed
	randStr := "123456" // 6-letter random string
	hashInput := fmt.Sprintf("%s/%s?%s#%s", randStr, methodName, paramString, secret)

	// 4. Calculate the SHA512 hash
	hasher := sha512.New()
	hasher.Write([]byte(hashInput))
	apiSig := hex.EncodeToString(hasher.Sum(nil))

	// 5. Construct the final API URL
	finalURL := fmt.Sprintf("https://codeforces.com/api/%s?%s&apiSig=%s%s", methodName, paramString, randStr, apiSig)

	// 6. Make the API call and check the status
	resp, err := http.Get(finalURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Status == "OK", nil
}

// Fetches all problems from the problemset
func fetchAllProblems() ([]CFProblem, error) {
	url := "https://codeforces.com/api/problemset.problems"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Problems []CFProblem `json:"problems"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Result.Problems, nil
}

// Fetches a user's submissions and returns a map of solved problem IDs
func fetchSolvedProblems(handle string) (map[string]bool, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s", handle)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result []CFSubmission `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	solved := make(map[string]bool)
	for _, sub := range result.Result {
		if sub.Verdict == "OK" {
			problemID := fmt.Sprintf("%d%s", sub.Problem.ContestID, sub.Problem.Index)
			solved[problemID] = true
		}
	}
	return solved, nil
}

func VerifyCodeforcesHandle(handle string) (bool, error) {
	if handle == "" {
		return false, errors.New("handle cannot be empty")
	}

	url := fmt.Sprintf("https://codeforces.com/api/user.info?handles=%s", handle)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// The API returns "OK" if the user exists and "FAILED" otherwise.
	if result.Status == "OK" {
		return true, nil
	}

	return false, nil
}

// --- NEW FUNCTION to get practice problems for a user ---
func GetPracticeProblems(handle string, rating int) ([]CFProblem, error) {
	allProblems, err := fetchAllProblems()
	if err != nil {
		return nil, err
	}

	solvedProblems, err := fetchSolvedProblems(handle)
	if err != nil {
		return nil, err
	}

	var practiceProblems []CFProblem
	for _, problem := range allProblems {
		problemID := fmt.Sprintf("%d%s", problem.ContestID, problem.Index)
		// Find unsolved problems that match the desired rating
		if !solvedProblems[problemID] && problem.Rating == rating {
			practiceProblems = append(practiceProblems, problem)
			// Limit to 10 problems for this example
			if len(practiceProblems) >= 10 {
				break
			}
		}
	}
	return practiceProblems, nil
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
