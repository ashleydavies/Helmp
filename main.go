package main

import "bytes"
import "strings"
import "text/template"
import "gopkg.in/yaml.v2"
import "github.com/gopherjs/gopherjs/js"
import "github.com/Masterminds/sprig"
import "honnef.co/go/js/dom"

func main() {
	js.Global.Get("template").Call("addEventListener", "input", func() {
		go update()
	})
	js.Global.Get("values").Call("addEventListener", "input", func() {
		go update()
	})
}

func update() {
	doc := dom.GetWindow().Document()
	templateText := doc.GetElementByID("template").(*dom.HTMLTextAreaElement).Value
	valuesText := doc.GetElementByID("values").(*dom.HTMLTextAreaElement).Value

	values := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(valuesText), &values)
	allThings := struct{
		Chart struct{
			Name string
			Version string
		}
		Release struct{
			Name string
			Service string
		}
		Values map[interface{}]interface{}
	}{
		Chart: struct{
			Name string
			Version string
		}{
			Name: "ChartName",
			Version: "ChartVersion",
		},
		Release: struct{
			Name string
			Service string
		}{
			Name: "ReleaseName",
			Service: "ReleaseService",
		},
		Values: values,
	}

	sprigMap := sprig.TxtFuncMap()
	sprigMap["toYaml"] = toYAML
	tmpl, err := template.New("user_template").Funcs(sprigMap).Parse(templateText)
	if err != nil {
		dom.GetWindow().Document().GetElementByID("output").SetInnerHTML(err.Error())
	}
	buffer := new(bytes.Buffer)
	err = tmpl.Execute(buffer, allThings)
	if err != nil {
		dom.GetWindow().Document().GetElementByID("output").SetInnerHTML(err.Error())
	}
	dom.GetWindow().Document().GetElementByID("output").SetInnerHTML(buffer.String())
}

// Stolen from Helm source
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
