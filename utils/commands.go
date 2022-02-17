package utils

import (
	"bytes"
	"context"
	"github.com/codeclysm/extract"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func runCommand(baseCmd string, workDir string, args []string) []byte {
	cmd := exec.Command(baseCmd, args...)
	cmd.Dir = workDir
	log.Printf("Running command: %v\n", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n%v\n", out, err)
	}
	log.Printf("Command result: %s", string(out))
	return out
}

func getCcoImageDigest(pullSecretFile string, imageUrl string) string {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "info", "--image-for", "cloud-credential-operator", imageUrl}
	log.Printf("Obtaining Cloud Credentials Operator image digest from image: %v\n", imageUrl)
	out := runCommand(baseCmd, "", args)

	return strings.TrimSuffix(string(out), "\n")
}

func findTarballs(outputDir string) []string {
	baseCmd := "find"
	args := []string{outputDir, "-name", "*.tar.*"}
	log.Printf("Looking up tarballs in : %v", outputDir)
	out := runCommand(baseCmd, "", args) //Must not switch dir.

	return strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
}

func Unarchive(outputDir, targetDir string) {
	tarballs := findTarballs(outputDir)
	for _, tarball := range tarballs {
		log.Printf("Extracting: %v", tarball)
		data, _ := ioutil.ReadFile(tarball)
		buffer := bytes.NewBuffer(data)
		extract.Gz(context.TODO(), buffer, targetDir, nil)
	}
	//TODO: handle errors
}

func ExtractTools(pullSecretFile string, outputDir string, imageUrl string) {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "extract", "--tools", imageUrl}
	log.Printf("Extracting tools from image: %v", imageUrl)
	runCommand(baseCmd, outputDir, args)
}

func ExtractCcoctl(pullSecretFile string, outputDir string, imageUrl string) {
	ccoImage := getCcoImageDigest(pullSecretFile, imageUrl)
	baseCmd := "oc"
	args := []string{"image", "-a", pullSecretFile, "extract", "--file", "/usr/bin/ccoctl", ccoImage}
	log.Printf("Extracting ccoctl from image: %v", ccoImage)
	runCommand(baseCmd, outputDir, args)

	//TODO: make sure ccoctl is executable.

}

func InstallCluster(installDir string, verbose bool) {
	baseCmd := "./openshift-install"
	args := []string{"create", "cluster"}
	if verbose {
		args = append(args, "-v")
	}
	log.Printf("Starting cluster installation.")
	runCommand(baseCmd, installDir, args)
}

func DestroyCluster(installDir string, verbose bool) {
	baseCmd := "./openshift-install"
	args := []string{"destroy", "cluster"}
	if verbose {
		args = append(args, "-v")
	}
	log.Printf("Destroying cluster.")
	runCommand(baseCmd, installDir, args)
}
