package main

import (
	"fmt"
	"github.com/RomanBednar/install-tools/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = "conf"

	// The environment variable prefix of all environment variables bound to command line flags.
	// For example, if the flag is --cloud, the environment variable will be INST_CLOUD
	envPrefix = "INST"

	secretsDir = "./secrets"
)

var (
	configPaths = []string{
		// First path here has the highest priority.
		"$HOME/.install-tools",
		"./config",
	}
)

func init() {
	cobra.OnInitialize(initializeConfig)

	rootCmd.PersistentFlags().StringP("action", "a", "create", "Action to perform. Valid values are: create, destroy.")
	viper.BindPFlag("action", rootCmd.PersistentFlags().Lookup("action"))

	rootCmd.PersistentFlags().StringP("cloud", "c", "aws", "Cloud to use for installation.")
	viper.BindPFlag("cloud", rootCmd.PersistentFlags().Lookup("cloud"))

	//TODO: add a scraper to resolve image url by version only (e.g. --image 4.10.0-rc.2)
	rootCmd.PersistentFlags().StringP("image", "i", "", "OpenShift image to use for installation. HINT: get full URL image at https://amd64.ocp.releases.ci.openshift.org/")
	viper.BindPFlag("image", rootCmd.PersistentFlags().Lookup("image"))

	rootCmd.PersistentFlags().StringP("cluster-name", "n", "mytestcluster-1", "Name of the cluster to create.")
	viper.BindPFlag("clustername", rootCmd.PersistentFlags().Lookup("cluster-name"))

	rootCmd.PersistentFlags().StringP("user-name", "u", "mytestuser-1", "Name of the user to create.")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("user-name"))

	rootCmd.PersistentFlags().StringP("output-dir", "o", "./_output", "Directory to write output files to.")
	viper.BindPFlag("outputdir", rootCmd.PersistentFlags().Lookup("output-dir"))

	rootCmd.PersistentFlags().StringP("cloud-region", "r", "us-east-1", "Cloud region to use for installation.")
	viper.BindPFlag("cloudregion", rootCmd.PersistentFlags().Lookup("cloud-region"))

	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "Dry run - only generate install-config.yaml and manifests, do not install cluster.")
	viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dry-run"))
}

func initializeConfig() {
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	viper.SetConfigName(defaultConfigFilename)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there is no a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Config file not found: %v\n", err)
		} else {
			log.Fatal(err)
		}
	}
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func validateFlags() {
	if viper.GetString("image") == "" {
		log.Fatalf("Image must be specified.")
	}
	if viper.GetString("cloud") == "" {
		log.Fatalf("Cloud must be specified.")
	}
}

var rootCmd = &cobra.Command{
	Use:   "install-tool",
	Short: "OpenShift install tool",
	Long:  `Simple tool for installing OpenShift on various clouds.`,
	Run: func(cmd *cobra.Command, args []string) {
		var c utils.Config
		if err := viper.Unmarshal(&c); err != nil {
			fmt.Printf("Error unmarshalling config file: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Running with configuration: %#v\n", c)
		Run(&c)
	},
}

func Run(conf *utils.Config) {
	//os.Exit(0)
	validateFlags()

	utils.MustDockerLogin(secretsDir, conf.Image)

	parser := utils.NewTemplateParser(conf)
	parser.ParseTemplate()

	utils.NewInstallDriver(conf).Run()

	if conf.DryRun {
		log.Printf("Done.")
		return
	}

	action := viper.GetString("action")
	switch action {
	case "create":
		utils.InstallCluster(conf.OutputDir, true)
	case "destroy":
		utils.DestroyCluster(conf.OutputDir, true)
	default:
		log.Fatalf("Unkown action: %v. Exiting.", action)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
