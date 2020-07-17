# Embedded GoTemplates
This package allows you to encode and bundle templates into your go binary without fear of escape sequences

### Example
1. `go get github.com/harryrford/template-embedded`
2. `template-embedded -in=./example -out=./example/generate.gen.go -package=example`
<details><summary>Generated Output</summary>
<p>

```go
package example

import (
	"encoding/base64"
	"text/template"
)

var embeddedTemplates = map[string]string{
	"example.tmpl": "cGFja2FnZSB7ey5QYWNrYWdlTmFtZX19DQoNCmltcG9ydCAiZm10Ig0KDQpmdW5jIG1haW4oKSB7DQoJZm10LlByaW50bG4oInt7Lk1lc3NhZ2V9fSIpDQp9",
}

// Parse parses declared templates
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
}
```

</p>
</details>