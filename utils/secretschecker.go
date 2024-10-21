package utils

import (
	"log"
	"net/url"
	"os"
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

func MustContainerEngineLogin(pullSecretFile, imageUrl, engine string) {
	registryDomain := getDomainFromURL(imageUrl)
	baseCmd := engine
	args := []string{"login", "--authfile", os.ExpandEnv(pullSecretFile), registryDomain}
	log.Printf("Verifying we can login with %v to: %v", engine, registryDomain)
	_, _, rc := runCommand(baseCmd, "", args...)

	if rc != 0 {
		panic("Could not login to registry: " + registryDomain)
	}
}
