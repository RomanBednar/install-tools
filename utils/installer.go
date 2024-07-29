package utils

import (
	"fmt"
	"log"
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
	case "vmware":
		fmt.Println("Driver is preparing vmWare installation.")
		d.vmwarePreparation()
	case "alibaba":
		fmt.Println("Driver is preparing Alibaba installation.")
		d.alibabaPreparation()
	case "azure":
		fmt.Println("Driver is preparing Azure installation.")
		d.azurePreparation()
	default:
		panic(fmt.Errorf("Unsupported cloud selected: %v\n", d.conf.Cloud))
	}

}

func (d *InstallDriver) awsPreparation() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	// Unarchive(d.conf.OutputDir, d.conf.OutputDir)
}

// For installing EFS Operator via Operator Hub refer to documentation provided there.
// Users have to create CredentialsRequest manually and let ccoctl create iam role - although similar this CredentialsRequest has nothing to do with the one created by the operator later.
// For --identity-provider-arn in ccoctl use existing identity provider that was used to create other roles by the installer.
func (d *InstallDriver) awsSTSPreparation() {
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	CreateInstallManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, "aws")
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	ExecuteCcoctl(d.conf.OutputDir, "aws", "us-east-1", d.conf.DryRun)
}

func (d *InstallDriver) gcpWIFPreparation() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	CreateInstallManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, "gcp")
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)

	//NOTE: for some reason the region for ccoctl binary does not match region in install-config.yaml
	ExecuteCcoctl(d.conf.OutputDir, "gcp", "us", d.conf.DryRun)
}

func (d *InstallDriver) vmwarePreparation() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	//Unarchive(d.conf.OutputDir, d.conf.OutputDir)
	// Start Bastion tunnel or scp install dir to bastion
	// TODO
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

func Run(conf *Config) {

	MustContainerEngineLogin(conf.PullSecretFile, conf.Image, conf.Engine)

	// This will start cluster installation/uninstallation.
	switch conf.Action {
	case "create":
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
