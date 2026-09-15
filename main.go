package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const neededArgs = 1

type Event struct {
	Type      string  `json:"type"`
	Repo      Repo    `json:"repo"`
	CreatedAt string  `json:"created_at"`
	Payload   Payload `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Payload struct {
	Action  string `json:"action"`
	RefType string `json:"ref_type"`
}

func main() {
	user, err := userName()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	bodyBytes, err := getActivity(user)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	events := make([]Event, 0)
	if err := json.Unmarshal(bodyBytes, &events); err != nil {
		fmt.Fprintf(os.Stderr, "Err: Cannot parse the JSON body")
		os.Exit(1)
	}

	fmt.Printf(eventSprinter(events))
	if len(events) == 0 {
		fmt.Println("No action has been made recently")
	}
}

func userName() (string, error) {
	if len(os.Args) != neededArgs+1 {
		var errMessage string
		if len(os.Args) > neededArgs+1 {
			errMessage = fmt.Sprintf("Too many arguments: needed %v, given %v\n", neededArgs, len(os.Args)-1)
		} else {
			errMessage = fmt.Sprintf("Too few arguments: needed %v, given %v\n", neededArgs, len(os.Args)-1)
		}
		return "", errors.New(errMessage)
	}
	return os.Args[neededArgs], nil
}

func getActivity(userName string) ([]byte, error) {
	url := urlMaker(userName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Failed to make the request!\n")
	}

	req.Header.Set("User-Agent", "github-user-activity")

	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		return nil, errors.New("Network error: could not reach GitHub!\n")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == 429 {
		remaining := resp.Header.Get("X-RateLimit-Remaining")
		resetStr := resp.Header.Get("X-RateLimit-Reset")
		if remaining == "0" && resetStr != "" {
			minutes, err := calNextReset(resetStr)
			if err != nil {
				return nil, err
			}
			if minutes <= 1 {
				return nil, errors.New("GitHub rate limit exceeded. Please wait a minute")
			} else {
				return nil, errors.New(fmt.Sprintf("GitHub rate limit exceeded. Please wait %v minutes and try again.\n", minutes))
			}
		}
		return nil, errors.New("Access forbidden by GitHub (possible rate limit)\n")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("An error occurred with the status code %v\n", resp.StatusCode))
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("An error occurred at reading the response body!\n")
	}
	return out, nil

}

func eventSprinter(events []Event) string {
	var out string
	for _, e := range events {
		switch e.Type {
		case "PushEvent":
			out += fmt.Sprintf("- Pushed commits to		%v\n", e.Repo.Name)
		case "IssuesEvent":
			out += fmt.Sprintf("- %v an issue in		%v\n", capitalise(e.Payload.Action), e.Repo.Name)
		case "WatchEvent":
			out += fmt.Sprintf("- Starred		%v\n", e.Repo.Name)
		case "PullRequestEvent":
			out += fmt.Sprintf("- %v a pull request in		%v\n", capitalise(e.Payload.Action), e.Repo.Name)
		case "CreateEvent":
			out += fmt.Sprintf("- Created a %v in		%v \n", e.Payload.RefType, e.Repo.Name)
		case "DeleteEvent":
			out += fmt.Sprintf("- Deleted a %v in		%v \n", e.Payload.RefType, e.Repo.Name)
		case "ForkEvent":
			out += fmt.Sprintf("- Forked		%v\n", e.Repo.Name)
		case "IssueCommentEvent":
			out += fmt.Sprintf("- Commented on an issue in		%v\n", e.Repo.Name)
		default:
			out += fmt.Sprintf("- Other activity in %v \n", e.Repo.Name)
		}
	}
	return out
}

func urlMaker(userName string) string {
	return fmt.Sprintf("https://api.github.com/users/%v/events", userName)
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func calNextReset(resetStr string) (int, error) {
	nextReset, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		return -1, errors.New("Cannot extract the reset time from the header\n")
	}
	now := time.Now().Unix()
	secondsLeft := nextReset - now
	return int((secondsLeft + 59) / 60), nil

}
