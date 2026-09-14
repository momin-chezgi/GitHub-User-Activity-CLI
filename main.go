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
	Type string `json:"type"`
	Repo Repo   `json:"repo"`
}

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
		fmt.Printf("%s: %s, (%s)\n", e.Type, e.Repo.Name, e.Repo.URL)
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
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("An error occurred with the status code %v\n", resp.StatusCode))
	}

	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return out, nil
}
