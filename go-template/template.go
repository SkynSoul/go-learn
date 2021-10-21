package gotmpl

import (
	"os"
	"text/template"
)

type Person struct {
	Name 	string
	Age		int
}

func Simple() {
	p := Person{"Li Lei", 18}
	t, err := template.New("simple").Parse(`Name: {{.Name}}, Age: {{.Age}}`)
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, p)
	if err != nil {
		panic(err)
	}
}
