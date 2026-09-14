package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
	Action string `json:"action"`
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

	for _, e := range events {
		switch e.Type {
		case "PushEvent":
			fmt.Printf("- Pushed a commit at (%v)\n", e.CreatedAt)
		case "IssuesEvent":
			fmt.Printf("- Opened/closed an issue at (%v)\n", e.CreatedAt)
		case "WatchEvent":
			fmt.Printf("- The user starred a repo/user at (%v)\n", e.CreatedAt)
		case "PullRequestEvent":
			fmt.Printf("- Opened/closed a pull request at (%v)\n", e.CreatedAt)
		case "CreateEvent":
			fmt.Printf("- Created a branch, tag, etc at (%v)\n", e.CreatedAt)
		case "DeleteEvent":
			fmt.Printf("- Deleted a branch, tag, etc at (%v)\n", e.CreatedAt)
		case "ForkEvent":
			fmt.Printf("- Forked a repo at (%v)\n", e.CreatedAt)
		case "IssueCommentEvent":
			fmt.Printf("- Commented on an issue at (%v)\n", e.CreatedAt)
		}
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
		err := errors.New(errMessage)
		return "", err
	}
	return os.Args[neededArgs], nil
}

func urlMaker(userName string) string {
	return fmt.Sprintf("https://api.github.com/users/%v/events", userName)
}

func getActivity(userName string) ([]byte, error) {
	url := urlMaker(userName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "github-user-activity")

	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("An error occurred with the status code %v\n", resp.StatusCode))
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return out, nil
}
