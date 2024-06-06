package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/codeclysm/extract"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const defaultFailedCode = 1

func runCommand(name string, workDir string, args ...string) (stdout string, stderr string, exitCode int) {
	log.Println("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)

	if exitCode != 0 {
		//TODO: maybe we should not panic here, but return error instead
		panic(fmt.Errorf("Command failed stderr=%v rc=%v", err, exitCode))
	}

	return
}

func getCcoImageDigest(pullSecretFile, outputDir, imageUrl string) string {
	// get absolute path of pullSecretFile
	file, err := filepath.Abs(pullSecretFile)
	if err != nil {
		panic(fmt.Sprintf("Could not resolve relative path to pull secret: %v", err))
	}

	baseCmd := "./oc"
	args := []string{"adm", "-a", file, "release", "info", "--image-for", "cloud-credential-operator", imageUrl}
	log.Printf("Obtaining Cloud Credentials Operator image digest from image: %v\n", imageUrl)
	out, _, _ := runCommand(baseCmd, outputDir, args...)

	return strings.TrimSuffix(out, "\n")
}

// Deprecated: findTarballs function is deprecated and will be removed in the future.
func findTarballs(outputDir string) []string {
	baseCmd := "find"
	args := []string{outputDir, "-name", "*.tar.*"}
	log.Printf("Looking up tarballs in : %v", outputDir)
	out, _, _ := runCommand(baseCmd, "", args...) //Must not switch dir.

	return strings.Split(strings.TrimSuffix(out, "\n"), "\n")
}

// Deprecated: Unarchive function is deprecated and will be removed in the future.
func Unarchive(outputDir, targetDir string) {
	log.Printf("Unarchiving tarballs from: %v to: %v", outputDir, targetDir)
	tarballs := findTarballs(outputDir)
	for _, tarball := range tarballs {
		log.Printf("Extracting: %v", tarball)
		data, err := os.ReadFile(tarball)
		if err != nil {
			log.Fatalf("Could not read tarball %s: %v", tarball, err)
		}
		buffer := bytes.NewBuffer(data)
		err = extract.Gz(context.TODO(), buffer, targetDir, nil)
		if err != nil {
			log.Fatalf("Could not extract tarball %s: %v", tarball, err)
		}
	}
}

// ExtractTools function extracts openshift-install and oc binaries from the image - this uses locally available oc binary
// which means it has to be run first and any consecutive commands should use the extracted oc binary.
func ExtractTools(pullSecretFile, outputDir, imageUrl string) {
	secret, err := filepath.Abs(os.ExpandEnv(pullSecretFile))
	if err != nil {
		panic(fmt.Sprintf("Could not resolve relative path to pull secret: %v", err))
	}
	baseCmd := "oc" //This has to be oc binary already present on the system because we don't have it extracted yet.

	//args := []string{"adm", "-a", secret, "release", "extract", "--tools", imageUrl}

	args := []string{"adm", "-a", secret, "release", "extract", "--command=openshift-install", imageUrl}
	log.Printf("Extracting openshift-install binary from image: %v", imageUrl)
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	args = []string{"adm", "-a", secret, "release", "extract", "--command=oc", imageUrl}
	log.Printf("Extracting oc binary from image: %v", imageUrl)
	_, _, _ = runCommand(baseCmd, outputDir, args...)
}

func ExtractCcoctl(pullSecretFile, outputDir, imageUrl string) {
	log.Printf("Extracting ccoctl from image: %v", imageUrl)
	// get absolute path of pullSecretFile
	file, err := filepath.Abs(pullSecretFile)
	if err != nil {
		panic(fmt.Sprintf("Could not resolve relative path to pull secret: %v", err))
	}

	ccoImage := getCcoImageDigest(file, outputDir, imageUrl)
	baseCmd := "./oc"
	args := []string{"image", "-a", file, "extract", "--file", "/usr/bin/ccoctl", "--confirm", ccoImage}
	log.Printf("Extracting ccoctl from CCO image digest: %v", ccoImage)
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "chmod"
	args = []string{"+x", "./ccoctl"}
	_, _, _ = runCommand(baseCmd, outputDir, args...)
}

func CreateCredentialRequestManifests(pullSecretFile, outputDir, imageUrl, region, cloud string) {
	// get absolute path of pullSecretFile
	file, err := filepath.Abs(pullSecretFile)
	if err != nil {
		panic(fmt.Sprintf("Could not resolve relative path to pull secret: %v", err))
	}

	log.Printf("Extracting manifests from image: %v", imageUrl)
	baseCmd := "./openshift-install"
	args := []string{"create", "manifests", "--log-level", "debug"}
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "awk"
	args = []string{"/infrastructureName:/{print $2}", "manifests/cluster-infrastructure-02-config.yml"}
	log.Println("Getting Infrastructure name")
	out, _, _ := runCommand(baseCmd, outputDir, args...)
	infrastructureName := strings.TrimSuffix(out, "\n")
	log.Printf("Infrastructure name found: %v", infrastructureName)

	baseCmd = "mkdir"
	args = []string{"creds", "cco-manifests"}
	log.Println("Creating creds directory.")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "./oc"
	args = []string{"adm", "-a", file, "release", "extract", "--credentials-requests", "--cloud", cloud, "--to", "./creds", imageUrl}
	log.Println("Extracting credential request")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "./ccoctl"
	args = []string{cloud, "create-ram-users", "--region", region, "--name", infrastructureName, "--credentials-requests-dir", "./creds", "--output-dir", "./cco-manifests"}
	log.Printf("Creating cloud credential manifests.")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	// Copy files to final manifests dir.
	path := filepath.Join(outputDir, "cco-manifests/manifests/*")
	files, _ := filepath.Glob(path)
	baseCmd = "cp"
	log.Printf("Copying cloud credential manifests to manifests dir.")
	for _, f := range files { //TODO: change this to one command
		args := []string{"-v", "-r", f, "./manifests"}
		_, _, _ = runCommand(baseCmd, outputDir, args...)
	}
}

func InstallCluster(installDir string, verbose bool) {
	baseCmd := "./openshift-install"
	args := []string{"create", "cluster"}
	if verbose {
		args = append(args, "--log-level", "debug")
	}
	log.Printf("Starting cluster installation.")
	_, _, _ = runCommand(baseCmd, installDir, args...)
	//TODO: this hides output from the progress - fix it
}

func DestroyCluster(installDir string, verbose bool) {
	baseCmd := "./openshift-install"
	args := []string{"destroy", "cluster"}
	if verbose {
		args = append(args, "--log-level", "debug")
	}
	log.Printf("Destroying cluster.")
	_, _, _ = runCommand(baseCmd, installDir, args...)
}
