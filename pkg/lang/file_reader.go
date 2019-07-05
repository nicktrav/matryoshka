package lang

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// fileReader controls how load() calls resolve and read other modules.
type fileReader interface {

	// Resolve parses the "name" part of load("name", "symbol") to a path. This
	// is not required to correspond to a true path on the filesystem, but should
	// be "absolute" within the semantics of this FileReader.
	//
	// fromPath will be empty when loading the root module passed to Load().
	Resolve(name, fromPath string) (path string, err error)

	// ReadFile reads the content of the file at the given path, which was
	// returned from Resolve().
	ReadFile(path string) ([]byte, error)
}

// localFile reads files from the local filesystem
type localFileReader struct {

	// root is the root directory with the files to parse
	root string
}

func (r *localFileReader) Resolve(name, fromPath string) (string, error) {
	if fromPath == "" {
		return name, nil
	}
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return "", fmt.Errorf("load(%q): invalid character in module name", name)
	}
	resolved := filepath.Join(r.root, filepath.FromSlash(path.Clean("/"+name)))
	return resolved, nil
}

func (r *localFileReader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
