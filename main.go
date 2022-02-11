package main

import (
	"github.com/spf13/viper"
	"os"
	"text/template"
)

func main() {
	os.Clearenv()
	template.New("test")
	viper.New()

}
