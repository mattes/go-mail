package mail

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chris-ramon/douceur/inliner"
	"github.com/yuin/goldmark"
	gmext "github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v2"
)

const (
	HTMLExtension     = "html"
	TextExtension     = "txt"
	YAMLExtension     = "yaml"
	ExampleFileSuffix = "example" // i.e. file.example.yaml
)

var defaultFuncMap = template.FuncMap{
	"markdown": parseMarkdown,
}

type Vars map[string]interface{}

type RenderOption func(*Templates)

type Templates struct {
	engine    *templateEngine
	files     Files
	inlineCSS bool // default true
}

func NewTemplates(files Files, opts ...RenderOption) (*Templates, error) {
	t := &Templates{
		engine:    newTemplateEngine(),
		files:     files,
		inlineCSS: true,
	}

	t.engine.funcMap = defaultFuncMap

	// apply options
	for _, fn := range opts {
		fn(t)
	}

	if err := t.engine.ParseFiles(files); err != nil {
		return nil, err
	}

	return t, nil
}

func OnMissingVar(mode MissingVar) RenderOption {
	return func(t *Templates) {
		t.engine.missingVar = mode
	}
}

func DisableInlineCSS() RenderOption {
	return func(t *Templates) {
		t.inlineCSS = false
	}
}

// Delims sets the delimiters
func SetDelims(left, right string) RenderOption {
	return func(t *Templates) {
		t.engine.leftDelim = left
		t.engine.rightDelim = right
	}
}

// Funcs adds the elements of the argument map to the template's function map.
func AddFuncs(funcMap template.FuncMap) RenderOption {
	return func(t *Templates) {
		for k, v := range funcMap {
			t.engine.funcMap[k] = v
		}
	}
}

// Func adds a function to the template's function map.
func AddFunc(name string, fn interface{}) RenderOption {
	return func(t *Templates) {
		t.engine.funcMap[name] = fn
	}
}

// RenderWithExample will render a template and use example data provided in name.example.yaml.
// It's goroutine-safe.
func (t *Templates) RenderWithExampleData(name string) ([]byte, error) {
	f, err := t.files.Open(exampleDataName(name))

	if err != nil && os.IsNotExist(err) {
		defer f.Close()

		// File doesn't exist, but maybe the template doesn't need vars?
		// Let's try to render and see what happens.
		out, errx := t.Render(name, nil)
		if errx != nil {
			return nil, multierr.Combine(err, errx)
		}
		return out, nil

	} else if err != nil {
		return nil, err
	}

	// file exists, let's read it and render the template
	defer f.Close()

	vars, err := readExampleData(f)
	if err != nil {
		return nil, err
	}

	return t.Render(name, vars)
}

// Render renders template with vars. It's goroutine-safe.
func (t *Templates) Render(name string, vars Vars) ([]byte, error) {
	body, err := t.engine.ExecuteTemplate(name, vars)
	if err != nil {
		return nil, err
	}

	// should we inline CSS if body is html?
	if t.inlineCSS && filepath.Ext(name) == "."+HTMLExtension {
		body, err = inlineCSS(body)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

func exampleDataName(name string) string {
	x := strings.TrimSuffix(name, filepath.Ext(name))
	x = x + "." + ExampleFileSuffix + "." + YAMLExtension
	return x
}

func readExampleData(r io.Reader) (Vars, error) {
	vars := make(Vars)
	d := yaml.NewDecoder(r)
	if err := d.Decode(&vars); err != nil {
		return nil, err
	}
	return vars, nil
}

// inlineCSS inlines CSS
func inlineCSS(html []byte) ([]byte, error) {
	out, err := inliner.Inline(string(html))
	if err != nil {
		return nil, err
	}
	return []byte(out), nil
}

// parseMarkdown is available as template command, it parses markdown to html.
func parseMarkdown(markdown string) (html template.HTML, err error) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(gmext.GFM),
		goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
	)

	err = md.Convert([]byte(markdown), &buf)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
