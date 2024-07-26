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

	rootCmd.PersistentFlags().StringP("pull-secret", "p", "", "Path to the pull secret file.")
	viper.BindPFlag("pullsecret", rootCmd.PersistentFlags().Lookup("pull-secret"))

	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "Dry run - only generate install-config.yaml and manifests, do not install cluster.")
	viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.PersistentFlags().StringP("config-path", "f", "", "Path to the configuration file (can be used in place of any flags).")
	viper.BindPFlag("configpath", rootCmd.PersistentFlags().Lookup("config-path"))

}

func initializeConfig() {

	configFilePath := viper.GetString("configpath")
	if configFilePath != "" {
		fmt.Printf("Using custom config path: %s\n", configFilePath)
		// Prepend the custom config path, so it has the highest priority in viper.
		configPaths = append([]string{configFilePath}, configPaths...)
	}

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	viper.SetConfigName(utils.DefaultConfigFilename)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there is no a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Config file not found: %v\n", err)
		} else {
			log.Fatal(err)
		}
	}
	viper.SetEnvPrefix(utils.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
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
		validateFlags()
		utils.Run(&c)
	},
}

func validateFlags() {
	if viper.GetString("image") == "" {
		log.Fatalf("Image must be specified.")
	}
	if viper.GetString("cloud") == "" {
		log.Fatalf("Cloud must be specified.")
	}
	if viper.GetBool("dryrun") && viper.GetString("action") != "create" {
		log.Fatalf("Dry run can only be used with create action.")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
