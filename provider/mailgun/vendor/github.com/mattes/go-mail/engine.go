package mail

import (
	"bytes"
	"fmt"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	textTemplate "text/template"
)

type MissingVar int

const (
	HaltOnError MissingVar = iota
	ContinueOnError
	ZeroValueOnError
)

// templateEngine combines html/template and text/template and calls functions
// based on extension of given template name.
type templateEngine struct {
	leftDelim, rightDelim string
	funcMap               map[string]interface{}
	missingVar            MissingVar

	htmlTemplate *htmlTemplate.Template
	textTemplate *textTemplate.Template
}

func newTemplateEngine() *templateEngine {
	return &templateEngine{
		leftDelim:  "{{",
		rightDelim: "}}",
		missingVar: HaltOnError,
	}
}

func (t *templateEngine) ExecuteTemplate(name string, vars Vars) ([]byte, error) {
	if !t.Exists(name) {
		return nil, fmt.Errorf("template '%v' not found", name)
	}

	var buf bytes.Buffer

	switch filepath.Ext(name) {
	case "." + HTMLExtension:
		if err := t.htmlTemplate.ExecuteTemplate(&buf, name, vars); err != nil {
			return nil, err
		}

	case "." + TextExtension:
		if err := t.textTemplate.ExecuteTemplate(&buf, name, vars); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (t *templateEngine) Exists(name string) bool {
	switch filepath.Ext(name) {
	case "." + HTMLExtension:
		for _, v := range t.htmlTemplate.Templates() {
			if name == v.Name() {
				return true
			}
		}

	case "." + TextExtension:
		for _, v := range t.textTemplate.Templates() {
			if name == v.Name() {
				return true
			}
		}
	}

	return false
}

func (t *templateEngine) ParseFiles(files Files) error {
	var err error

	t.textTemplate, err = parseTextFiles(files, t.leftDelim, t.rightDelim, t.funcMap)
	if err != nil {
		return err
	}

	t.htmlTemplate, err = parseHTMLFiles(files, t.leftDelim, t.rightDelim, t.funcMap)
	if err != nil {
		return err
	}

	switch t.missingVar {
	case ContinueOnError:
		t.textTemplate = t.textTemplate.Option("missingkey=default")
		t.htmlTemplate = t.htmlTemplate.Option("missingkey=default")

	case ZeroValueOnError:
		t.textTemplate = t.textTemplate.Option("missingkey=zero")
		t.htmlTemplate = t.htmlTemplate.Option("missingkey=zero")

	default:
		fallthrough
	case HaltOnError:
		t.textTemplate = t.textTemplate.Option("missingkey=error")
		t.htmlTemplate = t.htmlTemplate.Option("missingkey=error")
	}

	return nil
}

// TODO: try to combine parseTextFiles and parseHTMLFiles using an interface{}

// parseTextFiles is a modified copy of:
// https://godoc.org/github.com/golang/go/src/html/template#Template.ParseFiles
func parseTextFiles(files Files, leftDelim, rightDelim string, funcMap textTemplate.FuncMap) (*textTemplate.Template, error) {
	var t *textTemplate.Template

	err := files.Walk(func(name string, body io.ReadCloser) error {
		defer body.Close()

		if !strings.HasSuffix(name, "."+TextExtension) {
			return nil
		}

		// read full body
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		// create new template if necessary
		if t == nil {
			t = textTemplate.New(name).Delims(leftDelim, rightDelim).Funcs(funcMap)
		}

		// parse template
		var tmpl *textTemplate.Template
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		_, err = tmpl.Parse(string(b))
		return err
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}

// parseHTMLFiles is a modified copy of:
// https://godoc.org/github.com/golang/go/src/html/template#Template.ParseFiles
func parseHTMLFiles(files Files, leftDelim, rightDelim string, funcMap htmlTemplate.FuncMap) (*htmlTemplate.Template, error) {
	var t *htmlTemplate.Template

	err := files.Walk(func(name string, body io.ReadCloser) error {
		defer body.Close()

		if !strings.HasSuffix(name, "."+HTMLExtension) {
			return nil
		}

		// read full body
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		// create new template if necessary
		if t == nil {
			t = htmlTemplate.New(name).Delims(leftDelim, rightDelim).Funcs(funcMap)
		}

		// parse template
		var tmpl *htmlTemplate.Template
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		_, err = tmpl.Parse(string(b))
		return err
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}
