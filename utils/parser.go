package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

//All configuration is loaded into this structure and then used to parse templates.
type Config struct {
	SshPublicKeyFile string
	PullSecretFile   string
	ClusterName      string
	Username         string
	Password         string
	OutputDir        string

	SshPublicKey string
	PullSecret   string
}

type TemplateParser struct {
	data              Config
	requestedCloud    string
	templateDir       string
	outputFile        string
	cloudTemplatesMap map[string]string
}

func NewTemplateParser(requestedCloud string, data Config) TemplateParser {
	templateParser := TemplateParser{}

	templateParser.requestedCloud = requestedCloud
	templateParser.data = data

	//Flip file paths to string.
	templateParser.data.SshPublicKey = templateParser.fileToString(data.SshPublicKeyFile)

	//Base directory for templates.
	templateParser.templateDir = "templates/"

	//Output file name.
	templateParser.outputFile = "install-config.yaml"

	//Mapping from argument to file.
	templateParser.cloudTemplatesMap = map[string]string{
		"aws": "aws_basic.tmpl",
		//TODO add more templates
	}

	return templateParser
}

func (t *TemplateParser) getTemplatePath(name string) string {
	return t.templateDir + t.cloudTemplatesMap[name]
}

func (t *TemplateParser) getTemplateName(name string) string {
	templateName := t.cloudTemplatesMap[name]
	if templateName == "" {
		panic(fmt.Errorf("Template not found for requested cloud: %v\nUse one of: %q", name, t.getSupportedClouds()))
	}

	return templateName
}

func (t *TemplateParser) getSupportedClouds() []string {
	var keys []string
	for k, _ := range t.cloudTemplatesMap {
		keys = append(keys, k)
	}
	return keys
}

func (t *TemplateParser) fileToString(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)

}

func (t *TemplateParser) ParseTemplate() {
	templatePath := t.getTemplatePath(t.requestedCloud)
	templateName := t.getTemplateName(t.requestedCloud)
	fmt.Printf("Using template: %v with data: %+v\n", templatePath, t.data)

	template := template.Must(template.New(templateName).ParseFiles(templatePath))

	err := os.Mkdir(t.data.OutputDir, 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(fmt.Errorf("Could not create output dir: %v Error: %v", t.data.OutputDir, err))
	}

	output := t.data.OutputDir + "/" + t.outputFile
	//TODO: check file presence, overwrite?

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = template.Execute(f, t.data)
	if err != nil {
		panic(err)
	}

}
