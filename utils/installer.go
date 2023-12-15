package utils

import "fmt"

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
		fmt.Println("Driver starting AWS install flow.")
		d.awsInstallFlow()
	case "aws-odf": //TODO: this could be a parameter instead
		fmt.Println("Driver starting AWS ODF install flow.")
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
		panic(fmt.Errorf("Unsupported cloud install flow selected: %v\n", d.conf.Cloud))
	}

}

func (d *InstallDriver) awsInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	Unarchive(d.conf.OutputDir, d.conf.OutputDir)
}

func (d *InstallDriver) vmwareInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	Unarchive(d.conf.OutputDir, d.conf.OutputDir)
	// Start Bastion tunnel or scp install dir to bastion
	// TODO
}

func (d *InstallDriver) alibabaInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	Unarchive(d.conf.OutputDir, d.conf.OutputDir)

	// Extract ccoctl tool
	ExtractCcoctl(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	CreateCredentialRequestManifests(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image, d.conf.CloudRegion, "alibabacloud")
}

func (d *InstallDriver) azureInstallFlow() {
	// Extract and unarchive tools from image
	ExtractTools(d.conf.PullSecretFile, d.conf.OutputDir, d.conf.Image)
	Unarchive(d.conf.OutputDir, d.conf.OutputDir)
}
