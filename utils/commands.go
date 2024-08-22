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

const (
	defaultFailedCode     = 1
	defaultGcpProject     = "openshift-gce-devel"
	defaultCredRequestDir = "./credRequests"
)

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

func mustBeSupportedCloud(cloud string) {
	// check if cloud provided is one of supported values
	supportedClouds := []string{"gcp", "aws"}
	var supported bool
	for _, c := range supportedClouds {
		if c == cloud {
			supported = true
			break
		}
	}
	if !supported {
		panic(fmt.Sprintf("Unsupported cloud selected: %v\n", cloud))
	}
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
	log.Printf("Extracting CCO image from release image: %v", imageUrl)
	// get absolute path of pullSecretFile
	file, err := filepath.Abs(pullSecretFile)
	if err != nil {
		panic(fmt.Sprintf("Could not resolve relative path to pull secret: %v", err))
	}

	ccoImage := getCcoImageDigest(file, outputDir, imageUrl)
	baseCmd := "./oc"
	args := []string{"image", "-a", file, "extract", "--file", "/usr/bin/ccoctl", "--confirm", ccoImage}
	log.Printf("Extracting ccoctl binary from CCO image digest: %v", ccoImage)
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "chmod"
	args = []string{"+x", "./ccoctl"}
	_, _, _ = runCommand(baseCmd, outputDir, args...)
}

func CreateInstallManifests(pullSecretFile, outputDir, imageUrl, cloud string) {
	mustBeSupportedCloud(cloud)

	// get absolute path of pullSecretFile
	file, err := filepath.Abs(pullSecretFile)
	if err != nil {
		panic(fmt.Sprintf("Could not resolve relative path to pull secret: %v", err))
	}

	log.Printf("Extracting manifests from image: %v", imageUrl)
	baseCmd := "./openshift-install"
	args := []string{"create", "manifests", "--log-level", "debug"}
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "mkdir"
	args = []string{defaultCredRequestDir}
	log.Println("Creating creds directory.")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	baseCmd = "./oc"
	args = []string{"adm", "-a", file, "release", "extract", "--credentials-requests", "--cloud", cloud, "--to", defaultCredRequestDir, imageUrl}
	log.Println("Extracting credential request")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

	//baseCmd = "cp"
	//// This assumes that `openshift-install create manifests` command defaults output dir to ./manifests.
	//args = []string{"-a", "./manifests/tls", "."}
	//log.Println("Copying bound service account signing key to manifests dir.")
	//_, _, _ = runCommand(baseCmd, outputDir, args...)

}

// ExecuteCcoctl must run after CreateInstallManifests and ExtractCcoctl
func ExecuteCcoctl(outputDir, cloud, region string, dryRun bool) {
	mustBeSupportedCloud(cloud)

	baseCmd := "awk"
	args := []string{"/infrastructureName:/{print $2}", "manifests/cluster-infrastructure-02-config.yml"}
	log.Println("Getting Infrastructure name")
	out, _, _ := runCommand(baseCmd, outputDir, args...)
	infrastructureName := strings.TrimSuffix(out, "\n")
	log.Printf("Infrastructure name found: %v", infrastructureName)

	baseCmd = "./ccoctl"
	// Omitting --output-dir flag to let ccoctl save manifests to ./manifests (default) - from there we don't have to move it.
	args = []string{cloud, "create-all", "--name", infrastructureName, "--region", region, "--credentials-requests-dir", defaultCredRequestDir}
	switch cloud {
	case "gcp":
		args = append(args, []string{"--project", defaultGcpProject}...)
	case "aws":
		args = append(args, "--create-private-s3-bucket")
	}

	if dryRun {
		log.Println("Dry run requested, skipping ccoctl command.")
		log.Printf("To execute ccoctl command manually run: %v %v", baseCmd, strings.Join(args, " "))
		return
	}

	log.Printf("Creating cloud credential manifests.")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

}

func checkVCenterReachable() {
	url := "https://vcenter.devqe.ibmc.devcluster.openshift.com/"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Warning: URL %s is not reachable. Error: %v\n", url, err)
		panic("VCenter is not reachable. Please check your VPN connection and try again.")
	}
	defer resp.Body.Close()

	fmt.Printf("URL %s is reachable\n", url)
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

/////////////

// Deprecated
func alibabaCreateCredRequestManifests(pullSecretFile, outputDir, imageUrl, region, cloud string) {
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
