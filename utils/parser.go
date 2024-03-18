package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"golang.org/x/term"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"text/template"
)

func userConfirm() bool {
	prompt := promptui.Select{
		Label: "Select[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

// All configuration is loaded into this structure and then used to parse templates.
type Config struct {
	Action, Cloud, ClusterName, UserName, OutputDir, CloudRegion, Image, VmwarePassword string
	SshPublicKeyFile, SshPublicKey, PullSecretFile, PullSecret, ImageTag, Engine        string
	DryRun                                                                              bool
}

type TemplateParser struct {
	data              Config
	requestedCloud    string
	templateDir       string
	outputFile        string
	cloudTemplatesMap map[string]string
}

var cloudTemplatesMap = map[string]string{
	"aws":     "aws_basic.tmpl",
	"aws-odf": "aws_odf.tmpl",
	"vmware":  "vmware_basic.tmpl",
	"alibaba": "alibaba_basic.tmpl",
	"azure":   "azure_basic.tmpl",
}

func NewTemplateParser(data *Config) TemplateParser {
	log.Printf("Creating TemplateParser for cloud: %v\n", data.Cloud)
	log.Printf("TemplateParser data: %v\n", data)
	templateParser := TemplateParser{}

	templateParser.requestedCloud = data.Cloud
	templateParser.data = *data

	//Flip file paths to string.
	templateParser.data.SshPublicKey = templateParser.fileToString(data.SshPublicKeyFile, false)
	templateParser.data.PullSecret = templateParser.fileToString(data.PullSecretFile, true)

	//Base directory for templates.
	templateParser.templateDir = "templates/"

	//Output file name.
	templateParser.outputFile = "install-config.yaml"

	//Mapping from argument to file.
	templateParser.cloudTemplatesMap = cloudTemplatesMap
	log.Printf("TemplateParser created with cloud: %v\n", templateParser.requestedCloud)
	return templateParser
}

func (t *TemplateParser) getTemplatePath(name string) string {
	return t.templateDir + t.cloudTemplatesMap[name]
}

func (t *TemplateParser) getTemplateName(name string) string {
	fmt.Printf("Searching for template with name: %#v", name)
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

func (t *TemplateParser) fileToString(file string, compact bool) string {
	log.Printf("Reading file: %v\n", file)
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if compact {
		log.Printf("Compacting json file: %v\n", file)
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, content); err != nil {
			log.Fatal(err)
		}
		return buffer.String()
	}

	return string(content)

}

func (t *TemplateParser) ParseTemplate() {
	templatePath := t.getTemplatePath(t.requestedCloud)
	templateName := t.getTemplateName(t.requestedCloud)
	log.Printf("Using template: %v with data: %+v\n", templatePath, t.data)

	template := template.Must(template.New(templateName).ParseFiles(templatePath))

	output := filepath.Join(t.data.OutputDir, t.outputFile)
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		log.Printf("Output file %v already exists, overwrite?\n", output)
		if !userConfirm() {
			log.Fatalf("Aborting.")
		}
	}

	if t.data.Cloud == "vmware" {
		log.Printf("Are you connected to TwinGate VPN?\n", output)
		if !userConfirm() {
			log.Fatalf("Aborting.")
		}
		password := passwordPrompt("Please enter password for vcenter (vcenter.devqe.ibmc.devcluster.openshift.com)")
		t.data.VmwarePassword = password
	}

	//TODO: this probably should not be here - move to main?
	fmt.Printf("Creating output dir: %v\n", t.data.OutputDir)
	err := os.MkdirAll(t.data.OutputDir, 0755)
	if os.IsNotExist(err) {
		panic(fmt.Errorf("Could not create output dir: %v Error: %v", t.data.OutputDir, err))
	}

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = template.Execute(f, t.data)
	if err != nil {
		panic(err)
	}

	//TODO: maybe the install config should be backed up? openshift-install will destroy it

}

func passwordPrompt(prompt string) string {
	fmt.Printf("%s: ", prompt)
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	fmt.Print("\n")
	return string(bytepw)

}
