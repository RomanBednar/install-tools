package utils

import (
	"log"
	"os/exec"
)

func ExtractTools(pullSecretFile string, outputDir string, imageUrl string) {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "extract", "--tools", imageUrl}
	cmd := exec.Command(baseCmd, args...)
	cmd.Dir = outputDir
	log.Printf("Extracting tools from image: %v\n", imageUrl)

	log.Printf("Running command: %v\n", cmd)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n%v\n", stdoutStderr, err)
	}
	log.Printf("Command result: %s\n", stdoutStderr)

	//TODO: extract tar files obtained.
}

func GetCcoImageDigest(pullSecretFile string, outputDir string, imageUrl string) string {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "info", "--image-for", "cloud-credential-operator", imageUrl}
	cmd := exec.Command(baseCmd, args...)
	cmd.Dir = outputDir
	log.Printf("Obtaining Cloud Credentials Operator image digest from image: %v\n", imageUrl)

	log.Printf("Running command: %v\n", cmd)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n%v\n", stdoutStderr, err)
	}
	log.Printf("Command result: %s\n", stdoutStderr)

	return ""
}
