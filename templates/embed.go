package templates

import (
	"embed"
)

//go:embed *.tmpl
var F embed.FS

// ReadFile reads and returns the content of the named file.
func ReadFile(name string) ([]byte, error) {
	return F.ReadFile(name)
}

//// GetTemplatePath returns the path to the template file.
//func GetFilePath(filename string) (string, error) {
//	dir, _ := f.ReadDir(".")
//	fmt.Printf("dir: %v\n", dir)
//
//	// Find the file
//	for _, file := range dir {
//		if file.Name() == filename {
//			info, error := file.Info()
//			if error != nil {
//				return "", error
//			}
//			fmt.Printf("INFO: %v\n", info.Sys())
//		}
//	}
//	return "asfdasdgasg", nil
//
//}

//// GetTemplatePath returns the path to the template file.
//func GetTemplatePath(filename string) (string, error) {
//	// Read the directory
//	dir, err := fs.ReadDir(f, ".")
//	if err != nil {
//		return "", err
//	}
//
//	// Find the file
//	for _, file := range dir {
//		if file.Name() == filename {
//			return filepath.Abs(file.Name())
//		}
//	}
//
//	return "", fmt.Errorf("file not found: %s", filename)
//}
