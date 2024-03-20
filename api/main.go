package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.Println("Starting server on :8080")
	http.HandleFunc("/save", saveInstallerConfig)
	http.HandleFunc("/action", runAction)
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %#v", r)
	fmt.Fprintln(w, "Hello, world!")
}

func runAction(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %#v", r)
	if r.Method != http.MethodPost {
		fmt.Errorf("method not allowed: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Received body: %#v\n", r.Body)

	var action struct {
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received action: %#v", action)

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Action received successfully")

}

func saveInstallerConfig(w http.ResponseWriter, r *http.Request) {

	type Config struct {
		Username         string `json:"username"`
		SshPublicKeyFile string `json:"sshPublicKeyFile"`
		PullSecretFile   string `json:"pullSecretFile"`
		OutputDir        string `json:"outputDir"`
		ClusterName      string `json:"clusterName"`
		Image            string `json:"image"`
		CloudRegion      string `json:"cloudRegion"`
		Cloud            string `json:"cloud"`
		DryRun           string `json:"dryRun"`
	}

	log.Printf("Received request to store installerConfig: %#v", r)
	if r.Method != http.MethodPost {
		fmt.Errorf("method not allowed: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Received body: %#v\n", r.Body)

	var installerConfig Config
	if err := json.NewDecoder(r.Body).Decode(&installerConfig); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Storing config: %#v", installerConfig)

	file, err := os.OpenFile("/tmp/conf.env", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Errorf("error opening conf.env file: %v", err)
		http.Error(w, fmt.Sprintf("Error opening conf.env file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	log.Printf("Saved config to file: %v", installerConfig.OutputDir)

	format := fmt.Sprintf(`
userName=%s
sshPublicKeyFile=%s
pullSecretFile=%s
outputDir=%s
clusterName=%s
image=%s
cloudRegion=%s
cloud=%s
dryRun=%s`,
		installerConfig.Username,
		installerConfig.SshPublicKeyFile,
		installerConfig.PullSecretFile,
		installerConfig.OutputDir,
		installerConfig.ClusterName,
		installerConfig.Image,
		installerConfig.CloudRegion,
		installerConfig.Cloud,
		installerConfig.DryRun,
	)

	format = strings.TrimSpace(format)
	format += "\n"

	// Write data to conf.env file
	if _, err := fmt.Fprintf(
		file,
		format,
	); err != nil {
		http.Error(w, fmt.Sprintf("Error writing to conf.env file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Config stored successfully")
}
