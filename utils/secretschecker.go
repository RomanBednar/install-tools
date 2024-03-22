package utils

import (
	"log"
	"net/url"
	"os"
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

// TODO: add podman support
//func CanPodmanLogin(pullSecretFile, imageUrl string) bool {
//	registryDomain := getDomainFromURL(imageUrl)
//	//configPath := filepath.Dir(pullSecretFile)
//	baseCmd := "podman"
//	args := []string{"login", "--authfile", pullSecretFile, registryDomain}
//	log.Printf("Verifying we can 'podman login' to: %v", registryDomain)
//	_, _, rc := runCommand(baseCmd, "", args...)
//
//	return rc == 0
//}

func MustDockerLogin(pullSecretFile, imageUrl string) bool {
	pullSecretDir := filepath.Dir(os.ExpandEnv(pullSecretFile))
	registryDomain := getDomainFromURL(imageUrl)
	baseCmd := "docker"
	args := []string{"--config", pullSecretDir, "login", registryDomain}
	log.Printf("Verifying we can login with docker to: %v", registryDomain)
	_, _, rc := runCommand(baseCmd, "", args...)

	return rc == 0
}
