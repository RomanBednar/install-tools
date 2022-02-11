package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var (
	userName = flag.String("username", "", "Username override")
)

var (
	configName = "config"
	defaults   = map[string]interface{}{
		"userName": "testuser-1",
	}
	configPaths = []string{
		// Lookup is done in the same order we add the paths.
		"~/.install-tools",
		"./config",
		//".",
	}
)

type Config struct {
	Username string
	Password string
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

}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	configureViper()
	//fmt.Printf("From viper: %v\n", viper.GetString("common.clustername"))
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Could not unmarshal config to struct: %v", err)
	}
	fmt.Printf("config from struct: %v", config)

}
