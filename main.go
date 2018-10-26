package main

import (
	"fmt"
	"log"
	"net/http"
)

func call(url string) (*http.Response, error) {
	client := http.DefaultClient
	resp, err := client.Get("http://localhost:8080/")
	if err != nil {
		return nil, fmt.Errorf("could not make http call to server: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] can't close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responds with not OK: %s", resp.Status)
	}

	return resp, nil
}

func main() {
	_, err := call()
	fmt.Println(err)
}
