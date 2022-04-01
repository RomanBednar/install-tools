package utils

import (
	"log"
	"net/url"
	"path/filepath"
	"strings"
)

func getDomainFromURL(imageURL string) string {
	//Parses imageURL - can be both with scheme/protocol or without
	u, err := url.Parse(imageURL)
	if err != nil {
		log.Fatal(err)
	}
	domain := ""
	if u.Scheme == "" {
		parts := strings.Split(u.String(), "/")
		domain = parts[0]
	} else {
		parts := strings.Split(u.Hostname(), ".")
		domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
	}
	return domain
}

func CanDockerLogin(pullSecretFile, imageUrl string) bool {
	domain := getDomainFromURL(imageUrl)
	configPath := filepath.Base(pullSecretFile)
	baseCmd := "docker"
	args := []string{"--config", configPath, "login", domain}
	log.Printf("Verifying docker can login to: %v", domain)
	_, _, rc := runCommand(baseCmd, "", args...)

	return rc == 0
}
