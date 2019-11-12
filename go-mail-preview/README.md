# go-mail-preview

Render templates with data in your browser. 

Data is used from a `my-template.example.yaml` file that lives alongside the actual template.

## Usage

```
go get github.com/mattes/go-mail/go-mail-preview

go-mail-preview -listen :8000 -templates /path/to/templates

open http://localhost:8000/preview/my-template.html
open http://localhost:8000/preview/my-template.txt
```

