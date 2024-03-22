package main

import (
	"encoding/json"
	"fmt"
	"github.com/RomanBednar/install-tools/utils"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	filePath   string
	fileEvents = make(chan string)
)

func main() {
	log.Println("Starting server on :8080")
	http.HandleFunc("/save", saveInstallerConfig)
	http.HandleFunc("/action", runAction)
	//http.HandleFunc("/log", fileHandler)
	http.HandleFunc("/hello", helloHandler)

	filePath = "/tmp/output/.openshift_install.log" //TODO: make this configurable

	//go fileWatcher()

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
		fmt.Errorf("error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received action: %#v", action)

	// Run installer
	var config utils.Config
	file, err := ini.Load("/tmp/conf.env")
	if err != nil {
		fmt.Printf("Failed to load config file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the INI file into the struct
	if err := file.MapTo(&config); err != nil {
		fmt.Printf("Failed to unmarshal config file: %v\n", err)
		os.Exit(1)
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

//func fileWatcher() {
//	fmt.Printf("Watching file: %s\n", filePath)
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer watcher.Close()
//
//	done := make(chan bool)
//
//	go func() {
//		for {
//			select {
//			case event := <-watcher.Events:
//				if event.Op&fsnotify.Write == fsnotify.Write {
//					fileEvents <- event.Name
//				}
//			case err := <-watcher.Errors:
//				fmt.Println("Error:", err)
//			}
//		}
//	}()
//
//	err = watcher.Add(filePath)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	<-done
//}
//
//func fileHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/event-stream")
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//
//	for {
//		select {
//		case event := <-fileEvents:
//			file, err := os.Open(event)
//			if err != nil {
//				http.Error(w, "Error opening file", http.StatusInternalServerError)
//				return
//			}
//			defer file.Close()
//
//			buf := make([]byte, 1024)
//			for {
//				n, err := file.Read(buf)
//				if err != nil {
//					break
//				}
//				w.Write(buf[:n])
//				w.(http.Flusher).Flush()
//				time.Sleep(1 * time.Second) // Adjust frequency of updates
//			}
//		}
//	}
//}
