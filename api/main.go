package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/RomanBednar/install-tools/utils"
	"gopkg.in/ini.v1"
)

var (
	filePath   string
	fileEvents = make(chan string)

	// Operation status tracking
	opMutex   sync.Mutex
	opStatus  = "idle" // idle, running, completed, error
	opError   string
	opMessage string
)

const (
	locationFilePath = "/tmp/.cache/config-location"
	defaultEngine    = "podman" // users are not allowed to change this now, requires container environment changes
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	log.Println("Starting server on :8080")
	http.HandleFunc("/save", corsMiddleware(saveInstallerConfig))
	http.HandleFunc("/action", corsMiddleware(runAction))
	http.HandleFunc("/log", corsMiddleware(logFileHandler))
	http.HandleFunc("/hello", corsMiddleware(helloHandler))
	http.HandleFunc("/status", corsMiddleware(statusHandler))
	http.HandleFunc("/check-dir", corsMiddleware(checkDirHandler))

	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %#v", r)
	fmt.Fprintln(w, "Hello, world!")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	opMutex.Lock()
	defer opMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  opStatus,
		"error":   opError,
		"message": opMessage,
	})
}

func checkDirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dirPath := r.URL.Query().Get("path")
	if dirPath == "" {
		http.Error(w, "path parameter required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Check if directory exists
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist - that's fine, it's "empty"
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists": false,
			"empty":  true,
		})
		return
	}

	if !info.IsDir() {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists": true,
			"empty":  false,
			"error":  "path is not a directory",
		})
		return
	}

	// Check if directory is empty
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists": true,
			"empty":  false,
			"error":  fmt.Sprintf("cannot read directory: %v", err),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"exists": true,
		"empty":  len(entries) == 0,
	})
}

func runAction(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %#v", r)
	if r.Method != http.MethodPost {
		fmt.Printf("method not allowed: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if already running
	opMutex.Lock()
	if opStatus == "running" {
		opMutex.Unlock()
		http.Error(w, "An operation is already running", http.StatusConflict)
		return
	}
	opMutex.Unlock()

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
		http.Error(w, fmt.Sprintf("Error reading config location file: %v", err), http.StatusInternalServerError)
		return
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
		http.Error(w, fmt.Sprintf("Failed to load config file: %v", err), http.StatusInternalServerError)
		return
	}

	// Unmarshal the INI file into the struct
	if err := file.MapTo(&config); err != nil {
		fmt.Printf("Failed to unmarshal config file: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to unmarshal config file: %v", err), http.StatusInternalServerError)
		return
	}
	// Add action to config
	config.Action = action.Action

	fmt.Printf("Running with configuration: %#v\n", config)

	// Set status to running
	opMutex.Lock()
	opStatus = "running"
	opError = ""
	opMessage = fmt.Sprintf("Running %s...", action.Action)
	opMutex.Unlock()

	// Run asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				opMutex.Lock()
				opStatus = "error"
				opError = fmt.Sprintf("%v", r)
				opMessage = "Operation failed"
				opMutex.Unlock()
				log.Printf("Operation panicked: %v", r)
			}
		}()

		utils.Run(&config)

		opMutex.Lock()
		opStatus = "completed"
		opError = ""
		opMessage = fmt.Sprintf("Operation %s completed successfully", action.Action)
		opMutex.Unlock()
		log.Printf("Operation %s completed successfully", action.Action)
	}()

	// Respond immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "running",
		"message": fmt.Sprintf("Operation %s started", action.Action),
	})
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
		Action           string `json:"action"`
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

	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("error opening conf.env file: %v", err)
		http.Error(w, fmt.Sprintf("Error opening conf.env file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	log.Printf("Saved config to file: %v", configFilePath)

	format := fmt.Sprintf(`userName=%s
sshPublicKeyFile=%s
pullSecretFile=%s
outputDir=%s
clusterName=%s
image=%s
cloudRegion=%s
cloud=%s
dryRun=%s
engine=%s
`,
		installerConfig.Username,
		installerConfig.SshPublicKeyFile,
		installerConfig.PullSecretFile,
		installerConfig.OutputDir,
		installerConfig.ClusterName,
		installerConfig.Image,
		installerConfig.CloudRegion,
		installerConfig.Cloud,
		installerConfig.DryRun,
		defaultEngine,
	)

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
	locationFile, err := os.OpenFile(locationFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

	// Reset operation status when new config is saved
	opMutex.Lock()
	opStatus = "idle"
	opError = ""
	opMessage = ""
	opMutex.Unlock()

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Config stored successfully",
	})
}

func logFileHandler(w http.ResponseWriter, r *http.Request) {

	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if a direct path is provided
	logPath := r.URL.Query().Get("path")
	if logPath == "" {
		// Load config file location from cache.
		locationFile, err := os.ReadFile(locationFilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading config location file: %v", err), http.StatusInternalServerError)
			return
		}
		// The config file path points to <outputDir>/conf.env, we need the parent dir
		configDir := filepath.Dir(string(locationFile))
		logPath = filepath.Join(configDir, ".openshift_install.log")
	}

	// Read the contents of the log file
	fmt.Printf("Reading log file from: %v\n", logPath)
	fileContents, err := os.ReadFile(logPath)
	if err != nil {
		// Return empty string if log doesn't exist yet (operation just started)
		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, "")
			return
		}
		http.Error(w, fmt.Sprintf("Error reading log file: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the file contents as the response body
	w.Header().Set("Content-Type", "text/plain")
	w.Write(fileContents)
}
