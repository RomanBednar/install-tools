package utils

import (
	"log"
	"net/url"
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

//func CanPodmanLogin(pullSecretFile, imageUrl string) bool {
//	registryDomain := getDomainFromURL(imageUrl)
//	//configPath := filepath.Dir(pullSecretFile)
//	baseCmd := "podman" //TODO: making this configurable is a problem because docker does not have --authfile flag
//	args := []string{"login", "--authfile", pullSecretFile, registryDomain}
//	log.Printf("Verifying we can 'podman login' to: %v", registryDomain)
//	_, _, rc := runCommand(baseCmd, "", args...)
//
//	return rc == 0
//}

func MustDockerLogin(pullSecretFile, imageUrl string) bool {
	registryDomain := getDomainFromURL(imageUrl)
	//configPath := filepath.Dir(pullSecretFile)
	baseCmd := "docker"
	args := []string{"login", registryDomain}
	log.Printf("Verifying we can login with docker to: %v", registryDomain)
	_, _, rc := runCommand(baseCmd, "", args...)

	return rc == 0
}
