package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RomanBednar/install-tools/templates"
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
	Action                  string `ini:"action"`
	Cloud                   string `ini:"cloud"`
	ClusterName             string `ini:"clusterName"`
	UserName                string `ini:"userName"`
	OutputDir               string `ini:"outputDir"`
	CloudRegion             string `ini:"cloudRegion"`
	Image                   string `ini:"image"`
	VSpherePassword         string `ini:"vSpherePassword"`
	VSphereBaseDomain       string `ini:"vSphereBaseDomain"`
	VSphereVCenterSubdomain string `ini:"vSphereVCenterSubdomain"`
	VSphereApiVIP           string `ini:"vSphereApiVIP"`
	VSphereIngressVIP       string `ini:"vSphereIngressVIP"`
	SshPublicKeyFile        string `ini:"sshPublicKeyFile"`
	SshPublicKey            string `ini:"sshPublicKey"`
	PullSecretFile          string `ini:"pullSecretFile"`
	PullSecret              string `ini:"pullSecret"`
	Engine                  string `ini:"engine"`
	DryRun                  bool   `ini:"dryRun"`
}

type TemplateParser struct {
	data              Config
	requestedCloud    string
	outputFile        string
	cloudTemplatesMap map[string]string
}

var cloudTemplatesMap = map[string]string{
	"aws":     "aws_basic.tmpl",
	"aws-sts": "aws_sts.tmpl",
	"aws-odf": "aws_odf.tmpl",
	"vsphere": "vsphere_basic.tmpl",
	"alibaba": "alibaba_basic.tmpl",
	"azure":   "azure_basic.tmpl",
	"gcp-wif": "gcp_wif.tmpl",
	"gcp":     "gcp_basic.tmpl",
}

func NewTemplateParser(data *Config) TemplateParser {
	log.Printf("Creating TemplateParser for cloud: %v\n", data.Cloud)
	log.Printf("TemplateParser data: %#v\n", data)
	templateParser := TemplateParser{}

	templateParser.requestedCloud = data.Cloud
	templateParser.data = *data

	//Flip file paths to string.
	templateParser.data.SshPublicKey = templateParser.fileToString(data.SshPublicKeyFile, false)
	templateParser.data.PullSecret = templateParser.fileToString(data.PullSecretFile, true)

	//Output file name.
	templateParser.outputFile = "install-config.yaml"

	//Mapping from argument to file.
	templateParser.cloudTemplatesMap = cloudTemplatesMap
	log.Printf("TemplateParser created with cloud: %v\n", templateParser.requestedCloud)
	return templateParser
}

func (t *TemplateParser) getTemplatePath(filename string) string {
	dir, error := templates.F.ReadDir(".")
	if error != nil {
		panic(fmt.Errorf("Template not found: %v\n", error))
	}

	// Find the file
	var absPath string
	for _, file := range dir {
		if file.Name() == filename {
			absPath, err := filepath.Abs(file.Name())
			if err != nil {
				panic(fmt.Errorf("Template not found: %v\n", err))
			}
			return absPath
		}
	}
	return absPath
}

func (t *TemplateParser) getTemplateName(name string) string {
	fmt.Printf("Searching template for: %#v\n", name)
	templateName, ok := t.cloudTemplatesMap[name]
	if templateName == "" || !ok {
		panic(fmt.Errorf("Template not found for requested cloud: %v\nUse one of: %q", name, t.getSupportedClouds()))
	}
	fmt.Printf("Found template: %s\n", templateName)

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
	expandedFilePath := os.ExpandEnv(file)
	content, err := os.ReadFile(expandedFilePath)
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
	templateFileName := t.getTemplateName(t.requestedCloud)

	log.Printf("Using template: %v with data: %+v\n", templateFileName, t.data)

	tmp := template.Must(template.New(templateFileName).ParseFS(templates.F, templateFileName))

	output := filepath.Join(t.data.OutputDir, t.outputFile)

	//TODO: This can work only for CLI - fix it.
	//if _, err := os.Stat(output); !os.IsNotExist(err) {
	//	log.Printf("Output file %v already exists, overwrite?\n", output)
	//	if !userConfirm() {
	//		log.Fatalf("Aborting.")
	//	}
	//}
	//
	//if t.data.Cloud == "vsphere" {
	//	log.Printf("Are you connected to TwinGate VPN?\n", output)
	//	if !userConfirm() {
	//		log.Fatalf("Aborting.")
	//	}
	//	password := passwordPrompt("Please enter password for vcenter (vcenter.devqe.ibmc.devcluster.openshift.com)")
	//	t.data.VSpherePassword = password
	//}

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

	err = tmp.Execute(f, t.data)
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
