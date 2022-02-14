package aws

import (
	"fmt"
	"github.com/RomanBednar/install-tools/utils"
	"os"
	"text/template"
)

/*func parse_template(data utils.Config, temp template.Template) {
	t := template.Must(template.New("todos").Parse("You have task named \"{{ .Name}}\" with description: \"{{ .Description}}\""))
	fmt.Printf("Aws parsed: %v", t)
}
*/

func ParseTemplate(data utils.Config) {
	//TODO: do not hardcode name here
	templateFile := utils.ResolveTemplate("aws")
	fmt.Printf("Using template: %v with data: %+v\n", templateFile, data)

	t := template.Must(template.New("aws_basic.tmpl").ParseFiles(templateFile))

	f, err := os.OpenFile("/tmp/123.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		panic(err)
	}

}
