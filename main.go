package main

import (
	"fmt"
	"github.com/RomanBednar/install-tools/utils"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

var (
	//Required flags.
	cloud = flag.String("cloud", "", "Required - cloud provider to use.") //TODO: load possible values from cloudTemplatesMap
	image = flag.String("image", "", "Required - URL of the desired image.")

	//Possible overrides of config values.
	userName       = flag.String("username", "", "userName override.")
	vmwarePassword = flag.String("vmwarepassword", "", "vmwarePassword override.") //TODO: handle passwords more securely
	outputDir      = flag.String("outputdir", "", "outputDir override.")
	cloudRegion    = flag.String("cloudregion", "", "cloudRegion override.")

	// Flow control flags.
	action = flag.String("action", "", "Action to perform. Choose from: [\"create\", \"destroy\"]")
	dryRun = flag.Bool("dryrun", false, "Prepare installation files only.")
)

var (
	configName = "config"
	defaults   = map[string]interface{}{
		"outputDir": "./output",
	}
	configPaths = []string{
		// First path in this slice has the highest priority.
		"$HOME/.install-tools",
		"./config",
	}
)

func configureViper() {
	viper.SetConfigName(configName)

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Can't read config file: %v \n", err))
	}
	viper.BindPFlags(flag.CommandLine)
}

func validateFlags() {
	if *image == "" {
		log.Fatalf("Image has to be specified. Use --image flag to set it.")
	}
	if *cloud == "" {
		log.Fatalf("Cloud provider has to be specified. Use --cloud flag to set it.")
	}

}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	validateFlags()

	configureViper()

	// Get config struct and unmarshal viper config to it.
	var config utils.Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Could not unmarshal config to struct: %v", err)
	}

	//log.Printf("Config from struct: %+v\n", config)
	//log.Printf("Config from viper: %v\n", viper.AllSettings())
	//log.Fatalf("err")

	if !utils.CanDockerLogin(config.PullSecretFile, *image) {
		log.Fatalf("Authentication failed for image repo: %v\nThis is most likely invalid or expired secret."+
			"Please check your secrets file: %v\n", *image, config.PullSecretFile)
	}

	//TODO: add possibility to resolve image url by version only (e.g. --image 4.10.0-rc.2)

	//TODO: flows need to move to cloud specific modules
	// Parse template.
	parser := utils.NewTemplateParser(*cloud, config)
	parser.ParseTemplate()

	utils.NewInstallDriver(*cloud, *image, config).Run()

	if *dryRun {
		log.Printf("Done.")
		return
	}

	switch *action {
	case "create":
		utils.InstallCluster(config.OutputDir, true)
	case "destroy":
		utils.DestroyCluster(config.OutputDir, true)
	default:
		log.Fatalf("Unkown action: %v. Exiting.", *action)
	}

}
