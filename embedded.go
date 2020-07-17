package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var (
		input       string
		output      string
		packageName string
	)
	flag.StringVar(&input, "in", "", `The template file or directory containing all templates wishing to embed`)
	flag.StringVar(&output, "out", "", `Where to output the generated embedded template file, stdout is default`)
	flag.StringVar(&packageName, "package", "main", `The package name for generated embedded template file`)
	flag.Parse()

	if input == "" {
		fmt.Println("Please specify a file or path to a template file")
		os.Exit(1)
	}

	var encoded = make(map[string]string)

	// Support inputting singular files or a directory
	if strings.HasSuffix(input, ".tmpl") || strings.HasSuffix(input, ".html") {
		f, err := os.Open(input)
		if err != nil {
			fmt.Println("Failed to read input template file:", err)
			os.Exit(1)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("Failed to read input template file:", err)
			os.Exit(1)
		}
		encoded[f.Name()] = base64.StdEncoding.EncodeToString(b)
	} else {
		files, err := ioutil.ReadDir(input)
		if err != nil {
			fmt.Println("Failed to read input directory:", err)
			os.Exit(1)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".tmpl") || strings.HasSuffix(file.Name(), ".html") {
				f, err := os.Open(input + "/" + file.Name())
				if err != nil {
					fmt.Println("Failed to read input template file:", err)
					os.Exit(1)
				}
				defer f.Close()
				b, err := ioutil.ReadAll(f)
				if err != nil {
					fmt.Println("Failed to read input template file:", err)
					os.Exit(1)
				}
				encoded[file.Name()] = base64.StdEncoding.EncodeToString(b)
			}
		}
	}

	// Write the encoded templates to an output file
	if !strings.HasSuffix(output, ".gen.go") {
		output = strings.TrimSuffix(output, ".go") + ".gen.go"
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	// build the generated go output file
	w.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	w.WriteString("import (\n" + `    "encoding/base64"` + "\n" + `    "text/template"` + "\n)\n\n")
	w.WriteString("var embeddedTemplates = map[string]string{\n")

	for name, encoding := range encoded {
		w.WriteString("    " + fmt.Sprintf(`"%s"`, name) + fmt.Sprintf(`: "%s",`, encoding) + "\n")
	}
	w.WriteString(`}` + "\n\n")
	w.WriteString(`// Parse parses declared templates
	func Parse(t *template.Template) (*template.Template, error) {
		for name, encoded := range embeddedTemplates {
			var tmpl *template.Template
			if t == nil {
				t = template.New(name)
			}
			if name == t.Name() {
				tmpl = t
			} else {
				tmpl = t.New(name)
			}
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				return nil, err
			}
			if _, err := tmpl.Parse(string(decoded)); err != nil {
				return nil, err
			}
		}
		return t, nil
	}`)
	if err := w.Flush(); err != nil {
		fmt.Println("Failed to generated templates:", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(output, buf.Bytes(), os.ModePerm); err != nil {
		fmt.Println("Failed to output generated templates:", err)
		os.Exit(1)
	}
}
