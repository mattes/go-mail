# go-mail

* Load template files from disk or embedded store.
* Preview rendered template files with example data in browser with [go-mail-preview tool](./go-mail-preview).
* Automatically [inlines CSS](https://github.com/aymerick/douceur).
* Powered by Go's [html/template](https://golang.org/pkg/html/template/) and [text/template](https://golang.org/pkg/text/template/) engine.
* Supports [markdown](https://github.com/yuin/goldmark) to HTML parsing inside template.
* Use [provider](./provider) to send an email.

## Usage

```go
import (
  "github.com/mattes/go-mail"
  "github.com/mattes/go-mail/provider/mailgun"
)

// load templates
tpl, err := mail.NewTemplates(mail.FilesFromLocalDir("./templates"))

// create mail envelope
m := mail.New()
m.Subject = "Advice to self"
m.To("Mattes", "mattes@example.com")
m.Template(tpl, "simple", mail.Vars{
  "Body": "no ice cream after midnight",
})

// send email with mailgun (or any other provider)
p, err := mailgun.New("mailgun-domain", "mailgun-key")
err = p.Send(m)
```

## Templates

### Structure

Emails can have a HTML and/or text body. Templates are recognized by
their file extension.

```
my-template.html         -> used for html body, processed with Go's html/template engine
my-template.txt          -> used for text body, processed with Go's text/template engine
my-template.example.yaml -> used for preview
```

### Embed templates into Go binary

To embed templates into your Go binary, you can use a tool like
[go.rice](https://github.com/GeertJohan/go.rice).

Install go.rice first:

```
go get github.com/GeertJohan/go.rice
go get github.com/GeertJohan/go.rice/rice
```

Update your go code:

```go
//go:generate rice embed-go

var MyTemplates = rice.MustFindBox("./path/to/templates")
```

Run `go generate` to generate the embedded file. See [files.go](./files.go) and
[files_test.go](./files_test.go) for an example.


### Nice templates

There is a couple of tested email templates available, please have a look at:

* https://github.com/wildbit/postmark-templates
* https://github.com/mailgun/transactional-email-templates
* https://github.com/leemunroe/responsive-html-email-template
* https://github.com/leemunroe/amp-email-templates

