package utils

import (
	"log"
	"os/exec"
	"strings"
)

func ExtractTools(pullSecretFile string, outputDir string, imageUrl string) {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "extract", "--tools", imageUrl}
	cmd := exec.Command(baseCmd, args...)
	cmd.Dir = outputDir
	log.Printf("Extracting tools from image: %v", imageUrl)

	log.Printf("Running command: %v", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %s\n%v", out, err)
	}
	log.Printf("Command result: %s", out)

	//TODO: extract tar files obtained.
}

func ExtractCcoctl(pullSecretFile string, outputDir string, imageUrl string) {
	ccoImage := getCcoImageDigest(pullSecretFile, imageUrl)
	baseCmd := "oc"
	args := []string{"image", "-a", pullSecretFile, "extract", "--file", "/usr/bin/ccoctl", ccoImage}
	cmd := exec.Command(baseCmd, args...)
	cmd.Dir = outputDir
	log.Printf("Extracting ccoctl from image: %v", ccoImage)

	log.Printf("Running command: %v", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %s\n%v", out, err)
	}
	log.Printf("Command result: %s", out)

	//TODO: make sure ccoctl is executable.

}

func getCcoImageDigest(pullSecretFile string, imageUrl string) string {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "info", "--image-for", "cloud-credential-operator", imageUrl}
	cmd := exec.Command(baseCmd, args...)

	log.Printf("Obtaining Cloud Credentials Operator image digest from image: %v\n", imageUrl)

	log.Printf("Running command: %v\n", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n%v\n", out, err)
	}
	log.Printf("Command result: %s", string(out))

	return strings.TrimSuffix(string(out), "\n")
}
