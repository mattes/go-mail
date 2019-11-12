package mail

import "testing"

func TestFilesFromRiceBox(t *testing.T) {
	tpl, err := NewTemplates(SimpleTemplate())
	if err != nil {
		t.Fatal(err)
	}

	if !tpl.engine.Exists("simple.html") {
		t.Fatal("template file not found")
	}
}
