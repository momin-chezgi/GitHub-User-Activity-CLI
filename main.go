package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const neededArgs = 1

type Event struct {
	Type      string  `json:"type"`
	Repo      Repo    `json:"repo"`
	Actor     Actor   `json:"actor"`
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

type Actor struct {
	Login string `json:"login"`
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
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf(eventSprinter(events))
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
		return nil, errors.New("Failed to send the request!\n")
	}

	req.Header.Set("User-Agent", "github-user-activity")

	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		return nil, errors.New("Connection failed!\n")
	}

	defer resp.Body.Close()

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
	for _, e := range events {
		switch e.Type {
		case "PushEvent":
			return fmt.Sprintf("- Pushed commits to		%v\n", e.Repo.Name)
		case "IssuesEvent":
			return fmt.Sprintf("- %v an issue in		%v\n", capitalise(e.Payload.Action), e.Repo.Name)
		case "WatchEvent":
			return fmt.Sprintf("- Starred		%v\n", e.Repo.Name)
		case "PullRequestEvent":
			return fmt.Sprintf("- %v a pull request in		%v\n", capitalise(e.Payload.Action), e.Repo.Name)
		case "CreateEvent":
			return fmt.Sprintf("- Created a %v in		%v \n", e.Payload.RefType, e.Repo.Name)
		case "DeleteEvent":
			return fmt.Sprintf("- Deleted a %v in		%v \n", e.Payload.RefType, e.Repo.Name)
		case "ForkEvent":
			return fmt.Sprintf("- Forked		%v\n", e.Repo.Name)
		case "IssueCommentEvent":
			return fmt.Sprintf("- Commented on an issue in		%v\n", e.Repo.Name)
		default:
			return fmt.Sprintf("- Other activity in %v \n", e.Repo.Name)
		}
	}
	return ""
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
