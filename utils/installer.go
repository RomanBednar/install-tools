package utils

import (
	"fmt"
	"log"
	"os"
)

type InstallDriver struct {
	conf *Config
}

func NewInstallDriver(conf *Config) *InstallDriver {
	installDriver := InstallDriver{conf}
	return &installDriver
}

func (d *InstallDriver) Run() {
	switch d.conf.Cloud {
	case "aws":
		fmt.Println("Driver is preparing AWS installation.")
		d.awsPreparation()
	case "aws-sts":
		fmt.Println("Driver is preparing AWS STS installation.")
		d.awsSTSPreparation()
	case "aws-odf": //TODO: this should be a parameter instead
		fmt.Println("Driver is preparing AWS ODF installation.")
		d.awsPreparation()
	case "gcp-wif":
		fmt.Println("Driver is preparing GCP WIF installation.")
		d.gcpWIFPreparation()
	case "gcp":
		fmt.Println("Driver is preparing GCP installation.")
		d.gcpPreparation()
	case "vsphere":
		fmt.Println("Driver is preparing vSphere installation.")
		d.vspherePreparation()
	case "alibaba":
		fmt.Println("Driver is preparing Alibaba installation.")
		d.alibabaPreparation()
	case "azure":
		fmt.Println("Driver is preparing Azure installation.")
		d.azurePreparation()
	case "azure-wi":
		fmt.Println("Driver is preparing Azure Workload Identity installation.")
		d.azureWIPreparation()
	default:
		panic(fmt.Errorf("Unsupported cloud selected: %v\n", d.conf.Cloud))
	}

}

func (d *InstallDriver) awsPreparation() {
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
}

// For installing EFS Operator via Operator Hub refer to documentation provided there.
// Users have to create CredentialsRequest manually and let ccoctl create iam role - although similar this CredentialsRequest has nothing to do with the one created by the operator later.
// For --identity-provider-arn in ccoctl use existing identity provider that was used to create other roles by the installer.
func (d *InstallDriver) awsSTSPreparation() {
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	CreateInstallManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, "aws")
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	ExecuteCcoctl(d.conf.OutputDir, "aws", "us-east-1", d.conf.ResourceGroup, d.conf.DryRun)
}

// Installing cluster on GCP requires a service account which is pruned every ~3 days.
func (d *InstallDriver) gcpWIFPreparation() {
	CreateGCPServiceAccount(d.conf.UserName, d.conf.OutputDir)
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	CreateInstallManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, "gcp")
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)

	//NOTE: for some reason the region for ccoctl binary does not match region in install-config.yaml
	ExecuteCcoctl(d.conf.OutputDir, "gcp", "us", d.conf.ResourceGroup, d.conf.DryRun)
}

// Installing cluster on GCP requires a service account which is pruned every ~3 days.
func (d *InstallDriver) gcpPreparation() {
	CreateGCPServiceAccount(d.conf.UserName, d.conf.OutputDir)
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
}

func (d *InstallDriver) vspherePreparation() {
	checkVCenterReachable()
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
}

// Deprecated
func (d *InstallDriver) alibabaPreparation() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	// Unarchive(d.conf.OutputDir, d.conf.OutputDir)

	// Extract ccoctl tool
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	alibabaCreateCredRequestManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, d.conf.CloudRegion, "alibabacloud")
}

func (d *InstallDriver) azurePreparation() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	// Unarchive(d.conf.OutputDir, d.conf.OutputDir)
}

func (d *InstallDriver) azureWIPreparation() {
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	CreateInstallManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, "azure")
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	ExecuteCcoctl(d.conf.OutputDir, "azure", "centralus", d.conf.ResourceGroup, d.conf.DryRun)
}

func Run(conf *Config) {

	MustContainerEngineLogin(conf.PullSecretFile, conf.Image, conf.Engine)

	// This will start cluster installation/uninstallation.
	switch conf.Action {
	case "create":
		fmt.Printf("Creating output dir: %v\n", conf.OutputDir)
		err := os.MkdirAll(conf.OutputDir, 0755)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("could not create output dir: %v Error: %v", conf.OutputDir, err))
		}

		// If installing workload identity cluster on Azure we need to pass sanitized resource group name on two places:
		// 1. ccoctl --name argument
		// 2. resourceGroupName in install-config
		// These have to match and not contain any special characters!!!
		if conf.Cloud == "azure-wi" {
			conf.ResourceGroup = SanitizeResourceGroupName(conf.ResourceGroup)
		}

		// This will create the install-config.yaml file and save to outputDir.
		parser := NewTemplateParser(conf)
		parser.ParseTemplate()

		// This will extract the tools from the image, unarchive them and save to outputDir.
		NewInstallDriver(conf).Run()

		// Stop here if dry run is requested.
		if conf.DryRun {
			log.Printf("Done.")
			return
		}

		// This will create the cluster.
		InstallCluster(conf.OutputDir, true)
	case "destroy":
		DestroyCluster(conf.OutputDir, true)
	default:
		log.Fatalf("Unkown action: %v. Exiting.", conf.Action)
	}
}
