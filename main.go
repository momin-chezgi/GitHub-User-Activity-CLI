package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const neededArgs = 1

func main() {
	user, err := userName()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	url := urlMaker(user)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	req.Header.Set("User-Agent", "github-user-activity")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	if connStat := checkStat(resp.StatusCode); connStat != nil {
		fmt.Fprintf(os.Stderr, connStat.Error())
		os.Exit(1)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf(string(bodyBytes))
}

func userName() (string, error) {
	if len(os.Args) != neededArgs+1 {
		var errMessage string
		if len(os.Args) > neededArgs+1 {
			errMessage = fmt.Sprintf("Too much arguments: needed %v, given %v\n", neededArgs, len(os.Args)-1)
		} else {
			errMessage = fmt.Sprintf("Too few arguments: needed %v, given %v\n", neededArgs, len(os.Args)-1)
		}
		err := errors.New(errMessage)
		return "", err
	}
	return os.Args[neededArgs], nil
}

func checkStat(sc int) error {
	sc /= 100
	switch sc {
	case 1:
		return errors.New("Informational error")
	case 3:
		return errors.New("Redirecting")
	case 4:
		return errors.New("Client error")
	case 5:
		return errors.New("Server error")
	default:
		return nil
	}
}

func urlMaker(userName string) string {
	return fmt.Sprintf("https://api.github.com/users/%v/events", userName)
}
