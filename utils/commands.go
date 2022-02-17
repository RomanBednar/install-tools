package utils

import (
	"bytes"
	"context"
	"github.com/codeclysm/extract"
	"io/ioutil"
	"log"
	"os/exec"
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
	return
}

//func runCommand(baseCmd string, workDir string, args []string) ([]byte, error) {
//	cmd := exec.Command(baseCmd, args...)
//	cmd.Dir = workDir
//	log.Printf("Running command: %v\n", cmd)
//	out, err := cmd.CombinedOutput()
//	if exitError, ok := err.(*exec.ExitError); ok {
//		return exitError.ExitCode()
//	}
//	if err != nil {
//		log.Fatalf("%s\n%v\n", out, err)
//	}
//	log.Printf("Command result: %s", string(out))
//	return out, nil
//}

func getCcoImageDigest(pullSecretFile string, imageUrl string) string {
	baseCmd := "oc"
	args := []string{"adm", "-a", pullSecretFile, "release", "info", "--image-for", "cloud-credential-operator", imageUrl}
	log.Printf("Obtaining Cloud Credentials Operator image digest from image: %v\n", imageUrl)
	out, _, _ := runCommand(baseCmd, "", args...)

	return strings.TrimSuffix(out, "\n")
}

func findTarballs(outputDir string) []string {
	baseCmd := "find"
	args := []string{outputDir, "-name", "*.tar.*"}
	log.Printf("Looking up tarballs in : %v", outputDir)
	out, _, _ := runCommand(baseCmd, "", args...) //Must not switch dir.

	return strings.Split(strings.TrimSuffix(out, "\n"), "\n")
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
	_, _, _ = runCommand(baseCmd, outputDir, args...)
}

func ExtractCcoctl(pullSecretFile string, outputDir string, imageUrl string) {
	ccoImage := getCcoImageDigest(pullSecretFile, imageUrl)
	baseCmd := "oc"
	args := []string{"image", "-a", pullSecretFile, "extract", "--file", "/usr/bin/ccoctl", ccoImage}
	log.Printf("Extracting ccoctl from image: %v", ccoImage)
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	//TODO: make sure ccoctl is executable.

}

func InstallCluster(installDir string, verbose bool) {
	baseCmd := "./openshift-install"
	args := []string{"create", "cluster"}
	if verbose {
		args = append(args, "--log-level", "debug")
	}
	log.Printf("Starting cluster installation.")
	_, _, _ = runCommand(baseCmd, installDir, args...)
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
