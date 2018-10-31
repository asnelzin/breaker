package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asnelzin/breaker"
)

func call(url string) (*http.Response, error) {
	client := http.DefaultClient
	resp, err := client.Get(url)
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

func handler(br breaker.Breaker) http.HandlerFunc {
	url := "http://localhost:8080/"
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := br.Call(func() (interface{}, error) {
			return call(url)
		})
		w.Write([]byte(fmt.Sprintf("breaker state: %d, error: %v", br.GetState(), err)))
	}
}

func main() {
	breaker := breaker.Breaker{
		Threshold:         3,
		InvocationTimeout: 5 * time.Second,
		ResetTimeout:      10 * time.Second,
	}

	http.HandleFunc("/", handler(breaker))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
