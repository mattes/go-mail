package mail

import (
	"fmt"
	nmail "net/mail"
)

type Mail struct {
	FromName, FromAddress string
	ToName, ToAddress     string
	Subject               string

	templates     *Templates
	templateName  string
	templateVars  Vars
	renderSuccess bool // TODO should we add mutex for it?

	HTML, Text []byte
}

func New() *Mail {
	return &Mail{}
}

func (m *Mail) From(name, address string) {
	m.FromName = name
	m.FromAddress = address
}

func (m *Mail) To(name, address string) {
	m.ToName = name
	m.ToAddress = address
}

func (m *Mail) FromStr() string {
	a := nmail.Address{Name: m.FromName, Address: m.FromAddress}
	return a.String()
}

func (m *Mail) ToStr() string {
	a := nmail.Address{Name: m.ToName, Address: m.ToAddress}
	return a.String()
}

// Template sets the template to be used for this email. Make sure to call
// mail.Render() afterwards.
// If both name.html and name.txt are available, both will be included in mail.
// If only name.html or name.txt are available, the one available will be used.
// If neither name.html nor name.txt are available, an error is returned when Rendered.
func (m *Mail) Template(t *Templates, name string, vars Vars) {
	if t == nil {
		panic("mail.Template: t cannot be nil")
	}

	m.templates = t
	m.templateName = name
	m.templateVars = vars
	m.renderSuccess = false
}

// Render renders HTML and/or Text portion of email with given Template.
func (m *Mail) Render() error {
	if m.renderSuccess {
		return nil
	}

	// make subject available as var in template
	m.templateVars["Subject"] = m.Subject

	htmlExists := false
	if m.templates.engine.Exists(m.templateName + "." + HTMLExtension) {
		htmlExists = true
	}

	textExists := false
	if m.templates.engine.Exists(m.templateName + "." + TextExtension) {
		textExists = true
	}

	if !htmlExists && !textExists {
		return fmt.Errorf("neither '%v.%v' nor '%v.%v' could be found", m.templateName, HTMLExtension, m.templateName, TextExtension)
	}

	var err error

	if htmlExists {
		m.HTML, err = m.templates.Render(m.templateName+"."+HTMLExtension, m.templateVars)
		if err != nil {
			return err
		}
	}

	if textExists {
		m.Text, err = m.templates.Render(m.templateName+"."+TextExtension, m.templateVars)
		if err != nil {
			return err
		}
	}

	// HTML and/or Text successfully rendered
	m.renderSuccess = true

	// remove reference to allow G
	m.templates = nil
	m.templateName = ""
	m.templateVars = nil
	return nil
}
