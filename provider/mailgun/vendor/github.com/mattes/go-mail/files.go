package mail

//go:generate rice embed-go

import (
	"io"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
)

// Simple is a simple default template ready to be used.
// See https://github.com/leemunroe/responsive-html-email-template
func SimpleTemplate() Files {
	return FilesFromRiceBox(rice.MustFindBox("templates"))
}

type WalkFunc func(name string, body io.ReadCloser) error

type Files interface {
	Open(name string) (io.ReadCloser, error)
	Walk(WalkFunc) error
}

func FilesFromRiceBox(box *rice.Box) Files {
	return &riceBoxFiles{box}
}

type riceBoxFiles struct {
	*rice.Box
}

func (r *riceBoxFiles) Open(name string) (io.ReadCloser, error) {
	return r.Box.Open(name)
}

func (r *riceBoxFiles) Walk(fn WalkFunc) error {
	return r.Box.Walk("", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := r.Box.Open(path)
		if err != nil {
			return err
		}

		return fn(filepath.Base(path), f)
	})
}

func FilesFromLocalDir(dir string) Files {
	return &localFiles{dir}
}

type localFiles struct {
	dir string
}

func (d *localFiles) Open(name string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(d.dir, name))
}

func (d *localFiles) Walk(fn WalkFunc) error {
	return filepath.Walk(d.dir, func(path string, info os.FileInfo, err error) error {
		// If there was a problem walking to the file or directory named by path,
		// the incoming error will describe the problem.
		// If an error is returned processing stops.
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		return fn(filepath.Base(path), f)
	})
}
