package utils

import (
	"fmt"
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
}

type TemplateParser struct {
	data              Config
	requestedCloud    string
	templateDir       string
	cloudTemplatesMap map[string]string
}

func NewTemplateParser(requestedCloud string, data Config) TemplateParser {
	templateParser := TemplateParser{}

	templateParser.requestedCloud = requestedCloud
	templateParser.data = data

	//Base directory for templates.
	templateParser.templateDir = "templates/"
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
		panic(fmt.Errorf("Template not found for requested cloud: %v\nUse one of: %+v", name, t.cloudTemplatesMap))
	}

	return templateName
}

func (t *TemplateParser) ParseTemplate() {
	templatePath := t.getTemplatePath(t.requestedCloud)
	templateName := t.getTemplateName(t.requestedCloud)
	fmt.Printf("Using template: %v with data: %+v\n", templatePath, t.data)

	template := template.Must(template.New(templateName).ParseFiles(templatePath))

	f, err := os.OpenFile(t.data.OutputDir, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = template.Execute(f, t.data)
	if err != nil {
		panic(err)
	}

}
