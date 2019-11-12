package mailgun

import (
	"os"
	"testing"

	"github.com/mattes/go-mail"
)

func TestMailgun(t *testing.T) {
	// load files and parse templates
	files := mail.FilesFromLocalDir("../../templates")
	tpl, err := mail.NewTemplates(files)
	if err != nil {
		t.Fatal(err)
	}

	m := mail.New()
	m.FromAddress = os.Getenv("EXAMPLE_MAIL_ADDRESS")
	m.ToAddress = os.Getenv("EXAMPLE_MAIL_ADDRESS")
	m.Subject = "Signup for my app"
	m.Template(tpl, "simple", mail.Vars{
		"Preheader":       "",
		"Body":            "Please sign up for my app!",
		"CallToAction":    "Signup",
		"CallToActionURL": "https://www.google.com",
		"Footer":          "",
		"OrgAddress":      "San Francisco",
	})

	// send email with mailgun
	mg, err := New(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	if err != nil {
		t.Fatal(err)
	}
	if err := mg.Send(m); err != nil {
		t.Fatal(err)
	}
}
