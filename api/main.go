package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting server on :8080")
	http.HandleFunc("/storeCredentials", storeCredentialsHandler)
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: %#v", r)
	fmt.Fprintln(w, "Hello, world!")
}

func storeCredentialsHandler(w http.ResponseWriter, r *http.Request) {

	type Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	log.Printf("Received request to store credentials: %#v", r)
	if r.Method != http.MethodPost {
		fmt.Errorf("method not allowed: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Received body: %#v\n", r.Body)

	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Storing credentials: %#v", credentials)

	// Open conf.env file in append mode
	file, err := os.OpenFile("/tmp/conf.env", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening conf.env file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Write username and password to conf.env file
	if _, err := fmt.Fprintf(file, "USERNAME=%s\nPASSWORD=%s\n", credentials.Username, credentials.Password); err != nil {
		http.Error(w, fmt.Sprintf("Error writing to conf.env file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Credentials stored successfully")
}
