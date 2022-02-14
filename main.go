package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

var (
	username = flag.String("username", "", "Username override")
)

var (
	configName = "config"
	defaults   = map[string]interface{}{
		"username": "default-user",
	}
	configPaths = []string{
		// First path in this slice has the highest priority.
		"$HOME/.install-tools",
		"./config",
	}
)

type Config struct {
	SshPublicKeyFile string
	PullSecretFile   string
	ClusterName      string
	Username         string
	Password         string
}

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

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	configureViper()

	fmt.Printf("From viper: %v\n", viper.GetString("username"))
	fmt.Printf("From viper: %v\n", viper.AllSettings())

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Could not unmarshal config to struct: %v", err)
	}
	fmt.Printf("config from struct: %v", config)

}
