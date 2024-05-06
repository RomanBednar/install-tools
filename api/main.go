package main

import (
	"encoding/json"
	"fmt"
	"github.com/RomanBednar/install-tools/utils"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	filePath   string
	fileEvents = make(chan string)
)

const (
	locationFilePath = "/tmp/.cache/config-location"
)

func main() {
	log.Println("Starting server on :8080")
	http.HandleFunc("/save", saveInstallerConfig)
	http.HandleFunc("/action", runAction)
	http.HandleFunc("/log", logFileHandler)
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
		fmt.Printf("method not allowed: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Received body: %#v\n", r.Body)

	var action struct {
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		fmt.Printf("error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received action: %#v", action)

	log.Printf("Loading config file location from: %v\n", locationFilePath)
	configFilePath, err := os.ReadFile(locationFilePath)
	if err != nil {
		fmt.Printf("error reading config location file: %v", err)
	}

	configFile, err := os.ReadFile(string(configFilePath))
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		http.Error(w, fmt.Sprintf("Error reading config file: %v", err), http.StatusInternalServerError)
		return
	}

	// Run installer
	var config utils.Config
	file, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Failed to load config file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the INI file into the struct
	if err := file.MapTo(&config); err != nil {
		fmt.Printf("Failed to unmarshal config file: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to unmarshal config file: %v", err), http.StatusInternalServerError)
	}
	// Add action to config
	config.Action = action.Action

	fmt.Printf("Running with configuration: %#v\n", config)

	utils.Run(&config)

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
		fmt.Printf("method not allowed: %v", r.Method)
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

	configFilePath := filepath.Join(installerConfig.OutputDir, "conf.env")
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0770); err != nil {
		fmt.Printf("error creating output directory: %v", err)
		http.Error(w, fmt.Sprintf("Error creating output directory: %v", err), http.StatusInternalServerError)
		return
	}

	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error opening conf.env file: %v", err)
		http.Error(w, fmt.Sprintf("Error opening conf.env file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	log.Printf("Saved config to file: %v", configFilePath)

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

	// Store the config location in a file
	if err := os.MkdirAll(filepath.Dir(locationFilePath), 0770); err != nil {
		fmt.Printf("error creating cache directory: %v", err)
		http.Error(w, fmt.Sprintf("Error creating cache directory: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Storing config location to: %v\n", locationFilePath)
	locationFile, err := os.OpenFile(locationFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error opening config-location file: %v", err)
		http.Error(w, fmt.Sprintf("Error opening config-location file: %v", err), http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(
		locationFile,
		configFilePath,
	); err != nil {
		fmt.Printf("error writing to config-location file: %v", err)
		http.Error(w, fmt.Sprintf("Error writing to conf.env file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Config stored successfully")
}

func logFileHandler(w http.ResponseWriter, r *http.Request) {

	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load config file location from cache.
	locationFile, err := os.ReadFile(locationFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config location file: %v", err), http.StatusInternalServerError)
		return
	}
	logFile := filepath.Join(string(locationFile), "/.openshift_install.log")

	// Read the contents of the log file
	fmt.Printf("Reading log file from: %v\n", logFile)
	fileContents, err := os.ReadFile(logFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading log file: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the file contents as the response body
	w.Header().Set("Content-Type", "text/plain")
	w.Write(fileContents)
}
