package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {

	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "v1.0.0")
		})

		mux.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
				return
			}
			var req struct {
				InputString string `json:"inputString"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			decodedBytes, err := base64.StdEncoding.DecodeString(req.InputString)
			if err != nil {
				http.Error(w, "Invalid base64 string", http.StatusBadRequest)
				return
			}
			resp := struct {
				OutputString string `json:"outputString"`
			}{OutputString: string(decodedBytes)}
			json.NewEncoder(w).Encode(resp)
		})

		mux.HandleFunc("/hard-op", func(w http.ResponseWriter, r *http.Request) {
			sleepTime := time.Duration(10+rand.Intn(11)) * time.Second
			time.Sleep(sleepTime)
			if rand.Intn(2) == 0 {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				fmt.Fprintf(w, "Success")
			}
		})

		server := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		log.Println("Starting server on :8080")

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	time.Sleep(time.Second)

	client := &http.Client{}

	resp, err := client.Get("http://localhost:8080/version")

	if err != nil {
		log.Fatalf("Error fetching version: %v", err)
	}

	version, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Version: %s\n", version)
	resp.Body.Close()

	reqBody := `{"inputString":"SGVsbG8gd29ybGQ="}`
	resp, err = client.Post("http://localhost:8080/decode", "application/json", strings.NewReader(reqBody))
	if err != nil {
		log.Fatalf("Error decoding string: %v", err)
	}

	decoded, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Decoded string: %s\n", decoded)
	resp.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/hard-op", nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Error performing hard-op: %v", err)
	}
	
	hardOpResult, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Hard-op result: %s\n", hardOpResult)
	resp.Body.Close()
}
