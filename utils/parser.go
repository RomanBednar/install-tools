package utils

var (
	templateDir = "templates/"

	//Mapping from argument to file
	CloudTemplatesMap = map[string]string{
		"aws": "aws_basic.tmpl",
	}
)

//All configuration is loaded into this structure and then used to parse templates
type Config struct {
	SshPublicKeyFile string
	PullSecretFile   string
	ClusterName      string
	Username         string
	Password         string
}

//type TemplateParser struct {
//
//}

func ResolveTemplate(name string) string {
	return templateDir + CloudTemplatesMap[name]
}
