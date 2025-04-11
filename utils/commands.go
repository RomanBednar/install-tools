package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/codeclysm/extract"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	defaultFailedCode     = 1
	defaultGcpProject     = "openshift-gce-devel"
	defaultCredRequestDir = "./credRequests"

	// Name of a resource group we have preconfigured in Azure, used by ccoctl to find the right DNS zone
	defaultAzureResourceGroup = "os4-common"
)

func runCommand(name string, workDir string, args ...string) (stdout string, stderr string, exitCode int) {
	log.Println("run command:", name, strings.Join(args, " "))
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = strings.TrimSpace(outbuf.String())
	stderr = strings.TrimSpace(errbuf.String())

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
		panic(fmt.Errorf("command failed stderr=%v rc=%v", err, exitCode))
	}

	return
}

func mustBeSupportedCloud(cloud string) {
	// check if cloud provided is one of supported values
	supportedClouds := []string{"gcp", "aws", "azure"}
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
func ExecuteCcoctl(outputDir, cloud, region, rgName string, dryRun bool) {
	mustBeSupportedCloud(cloud)

	baseCmd := "./ccoctl"
	// Omitting --output-dir flag to let ccoctl save manifests to ./manifests (default) - from there we don't have to move it.
	args := []string{cloud, "create-all", "--name", rgName, "--region", region, "--credentials-requests-dir", defaultCredRequestDir}
	switch cloud {
	case "gcp":
		args = append(args, []string{"--project", defaultGcpProject}...)
	case "aws":
		args = append(args, "--create-private-s3-bucket")
	case "azure":
		azureAccount := getAzureCredentials()
		args = append(args, "--subscription-id", azureAccount.ID, "--dnszone-resource-group-name", defaultAzureResourceGroup, "--tenant-id", azureAccount.TenantID)
	}

	if dryRun {
		log.Println("Dry run requested, skipping ccoctl command.")
		log.Printf("To execute ccoctl command manually run: %v %v", baseCmd, strings.Join(args, " "))
		return
	}

	log.Printf("Creating cloud credential manifests.")
	_, _, _ = runCommand(baseCmd, outputDir, args...)

}

func CreateGCPServiceAccount(userName, outputDir string) {
	mustGcloudAuth()
	serviceAccountName := fmt.Sprintf("%s-development", userName)
	outputCredentialsFile := filepath.Join(outputDir, "gcp-service-account.json")
	baseCmd := "gcloud"
	var serviceAccountEmail string

	// First check if the account already exists
	args := []string{"iam", "service-accounts", "list", "--filter", fmt.Sprintf("displayName:%s", serviceAccountName), "--format", "value(displayName)"}
	output, _, _ := runCommand(baseCmd, "", args...)
	log.Printf("service account found: %#v needed: %#v", output, serviceAccountName)
	if output != serviceAccountName {
		// Create the service account
		log.Printf("Creating service account %s", serviceAccountName)
		args = []string{"iam", "service-accounts", "create", serviceAccountName, "--display-name", serviceAccountName}
		runCommand(baseCmd, "", args...)

		// Get service account email
		args = []string{"iam", "service-accounts", "list", "--filter", fmt.Sprintf("displayName:%s", serviceAccountName), "--format", "value(email)"}
		serviceAccountEmail, _, _ = runCommand(baseCmd, "", args...)
		if serviceAccountEmail == "" {
			log.Fatalf("Could not get service account email for %s", serviceAccountName)
			return
		}

		// Get service account project ID
		args = []string{"iam", "service-accounts", "list", "--filter", fmt.Sprintf("displayName:%s", serviceAccountName), "--format", "value(projectId)"}
		projectID, _, _ := runCommand(baseCmd, "", args...)
		if projectID == "" {
			log.Fatalf("Could not get project ID for %s", serviceAccountName)
			return
		}

		// Define permissions
		roles := []string{
			"roles/compute.admin",
			"roles/iam.securityAdmin",
			"roles/iam.serviceAccountAdmin",
			"roles/iam.serviceAccountKeyAdmin",
			"roles/iam.serviceAccountUser",
			"roles/storage.admin",
			"roles/dns.admin",
			"roles/compute.loadBalancerAdmin",
			"roles/iam.roleViewer",
			"roles/iam.workloadIdentityPoolAdmin",
		}

		//TODO: sometimes gcloud fails here with "There were concurrent policy changes. Please retry the whole read-modify-write with exponential backoff."
		//TODO: IAM commands should have a retry and backoff - this is a known issue with gcloud and is caused by some request limit per second which role creation easily exceeds.
		for _, role := range roles {
			args = []string{"projects", "add-iam-policy-binding", projectID, "--member", "serviceAccount:" + serviceAccountEmail, "--role", role, "--condition", "None"}
			runCommand(baseCmd, "", args...)
			time.Sleep(3 * time.Second) //TODO: fix this after exponential backoff is implemented
		}

		// Create service account key
		args = []string{"iam", "service-accounts", "keys", "create", outputCredentialsFile, "--iam-account", serviceAccountEmail}
		runCommand(baseCmd, "", args...)

	} else {
		// If storage account exists we would have to inspect the keys, save them as a file, make sure they're valid and what not - too complicated, it's easier to just recreate.
		args = []string{"iam", "service-accounts", "list", "--filter", fmt.Sprintf("displayName:%s", serviceAccountName), "--format", "value(email)"}
		serviceAccountEmail, _, _ = runCommand(baseCmd, "", args...)
		log.Printf("Service account %v already exists, please remove it and start again.", serviceAccountName)
		log.Printf("HINT: To remove service account run: gcloud iam service-accounts delete %s", serviceAccountEmail)
		panic("Installation aborted.")
	}

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", outputCredentialsFile); err != nil {
		panic(fmt.Sprintf("could not set GOOGLE_APPLICATION_CREDENTIALS env var: %v", err))
	}
	log.Printf("GOOGLE_APPLICATION_CREDENTIALS environment variable has been set to %s", outputCredentialsFile)
}

func mustGcloudAuth() {
	baseCmd := "gcloud"
	args := []string{"auth", "list", "--format", "json"}
	stdout, _, _ := runCommand(baseCmd, "", args...)

	var authList []map[string]string
	err := json.Unmarshal([]byte(stdout), &authList)
	if err != nil {
		panic(fmt.Sprintf("Error parsing gcloud auth list output: %v", err))
	}

	for _, auth := range authList {
		if auth["status"] == "ACTIVE" {
			return
		}
	}

	panic("Not logged in to gcloud. Please run 'gcloud auth login' first.")
}

func getAzureCredentials() azureAccountType {
	baseCmd := "az"
	args := []string{"account", "show", "-o", "json"}
	stdout, stderr, code := runCommand(baseCmd, "", args...)
	if code != 0 {
		panic(fmt.Sprintf("Error running \"az account show\", make sure to first log in with \"az login\": %s", stderr))
	}
	var azureAccount azureAccountType
	err := json.Unmarshal([]byte(stdout), &azureAccount)
	if err != nil {
		panic(fmt.Sprintf("Error parsing az account show output: %v", err))
	}
	return azureAccount
}

//func getInfrastructureName(dir string, sanitize bool) string {
//	baseCmd := "awk"
//	args := []string{"/infrastructureName:/{print $2}", "manifests/cluster-infrastructure-02-config.yml"}
//	log.Println("Getting Infrastructure name")
//	out, _, _ := runCommand(baseCmd, dir, args...)
//	infrastructureName := strings.TrimSuffix(out, "\n")
//	if sanitize {
//		// When passing a --name to ccoctl for Azure, the tool uses it for storage account name and has to be sanitized.
//		infrastructureName = sanitizeResourceGroupName(infrastructureName)
//	}
//	log.Printf("Infrastructure name found: %v", infrastructureName)
//
//	return infrastructureName
//}

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

type azureAccountType struct {
	EnvironmentName     string `json:"environmentName"`
	HomeTenantID        string `json:"homeTenantId"`
	ID                  string `json:"id"` //This is the same as Subscription ID
	IsDefault           bool   `json:"isDefault"`
	ManagedByTenants    []any  `json:"managedByTenants"`
	Name                string `json:"name"`
	State               string `json:"state"`
	TenantDefaultDomain string `json:"tenantDefaultDomain"`
	TenantDisplayName   string `json:"tenantDisplayName"`
	TenantID            string `json:"tenantId"`
	User                struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"user"`
}

// sanitizeResourceGroupName ensures the string meets Azure storage account requirements:
// - Between 3 and 24 characters
// - Only lowercase letters and numbers
// - Trims newline characters
func SanitizeResourceGroupName(name string) string {
	// First trim the newline character
	name = strings.TrimSuffix(name, "\n")

	// Convert to lowercase
	name = strings.ToLower(name)

	// Keep only lowercase letters and numbers
	var result strings.Builder
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			result.WriteRune(char)
		}
	}

	// Truncate to 24 characters if longer
	sanitized := result.String()
	if len(sanitized) > 24 {
		sanitized = sanitized[:24]
	}

	// Ensure at least 3 characters
	// If not enough valid characters, append placeholder digits
	for len(sanitized) < 3 {
		sanitized += "0"
	}

	return sanitized
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
