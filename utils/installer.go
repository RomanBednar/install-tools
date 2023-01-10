package utils

import "fmt"

type InstallDriver struct {
	cloud    string
	imageUrl string
	config   Config
}

func NewInstallDriver(cloudName string, imageUrl string, config Config) *InstallDriver {
	installDriver := InstallDriver{cloud: cloudName, imageUrl: imageUrl, config: config}
	return &installDriver
}

func (d *InstallDriver) Run() {
	switch d.cloud {
	case "aws":
		fmt.Println("Driver starting AWS install flow.")
		d.awsInstallFlow()
	case "vmware":
		fmt.Println("Driver starting vmWare install flow.")
		d.vmwareInstallFlow()
	case "alibaba":
		fmt.Println("Driver starting Alibaba install flow.")
		d.alibabaInstallFlow()
	case "azure":
		fmt.Println("Driver starting Azure install flow.")
		d.azureInstallFlow()
	default:
		panic(fmt.Errorf("Unsupported cloud install flow selected: %v\n", d.cloud))
	}

}

func (d *InstallDriver) awsInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.config.PullSecretFile, d.config.OutputDir, d.imageUrl)
	Unarchive(d.config.OutputDir, d.config.OutputDir)
}

func (d *InstallDriver) vmwareInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.config.PullSecretFile, d.config.OutputDir, d.imageUrl)
	Unarchive(d.config.OutputDir, d.config.OutputDir)
	// Start Bastion tunnel or scp install dir to bastion
	// TODO
}

func (d *InstallDriver) alibabaInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.config.PullSecretFile, d.config.OutputDir, d.imageUrl)
	Unarchive(d.config.OutputDir, d.config.OutputDir)

	// Extract ccoctl tool
	ExtractCcoctl(d.config.PullSecretFile, d.config.OutputDir, d.imageUrl)
	CreateCredentialRequestManifests(d.config.PullSecretFile, d.config.OutputDir, d.imageUrl, d.config.CloudRegion, "alibabacloud")
}

func (d *InstallDriver) azureInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.config.PullSecretFile, d.config.OutputDir, d.imageUrl)
	Unarchive(d.config.OutputDir, d.config.OutputDir)
}
