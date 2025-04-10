package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RomanBednar/install-tools/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initializeConfig)
	rootCmd.PersistentFlags().StringP("action", "a", "create", "Action to perform. Valid values are: create, destroy.")
	viper.BindPFlag("action", rootCmd.PersistentFlags().Lookup("action"))

	rootCmd.PersistentFlags().StringP("cloud", "c", "aws", fmt.Sprintf("Cloud to use for installation. Valid values are: %v", strings.Join(utils.GetCloudKeys(), ", ")))
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

	// dump configuration command
	rootCmd.PersistentFlags().BoolP("dump-config", "D", false, "Dump the configuration to stdout and exit.")
	viper.BindPFlag("dumpconfig", rootCmd.PersistentFlags().Lookup("dump-config"))

}

func initializeConfig() {

	configPaths := utils.ConfigPaths
	// If custom config path is used prepend it, so it has the highest priority in viper.
	configFilePath := viper.GetString("configpath")
	if configFilePath != "" {
		fmt.Printf("Using custom config path: %s\n", configFilePath)
		configPaths = append([]string{configFilePath}, utils.ConfigPaths...)
	}
	for _, path := range configPaths {
		viper.AddConfigPath(path)
		fmt.Printf("Added config path to viper: %s\n", path)
	}

	viper.SetConfigName(utils.DefaultConfigFilename)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there is no config file
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Printf("WARNING: Config file not found: %v\n", err)
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

		validateFlags()
		if viper.GetBool("dumpconfig") {
			fmt.Printf("Running with configuration: %#v\n", c)
			os.Exit(0)
		}

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
