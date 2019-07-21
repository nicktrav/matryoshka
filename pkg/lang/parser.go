package lang

import (
	"fmt"
	oslib "os"
	"path/filepath"

	"go.starlark.net/starlark"
)

const (
	main  = "main.dep"
	shell = "shell"
	dep   = "dep"
	os    = "os"
)

var (
	shellBuiltin = starlark.NewBuiltin(shell, FnShell)
	depBuiltin   = starlark.NewBuiltin(dep, FnDep)
	osBuiltin    = starlark.NewBuiltin(os, FnOs)

	defaultModules = starlark.StringDict{
		shell: shellBuiltin,
		dep:   depBuiltin,
		os:    osBuiltin,
	}
)

// Parser executes a Starlark program by parsing one or more Starlark files.
type Parser interface {

	// Run executes the parser.
	Run() error

	// Deps is a slice of pointers to parsed Deps.
	Deps() []*Dep
}

// NewParser returns a new Parser for the given root directory.
func NewParser(root string) Parser {
	return cachedParser{
		root:          root,
		reader:        &localFileReader{root},
		cache:         make(map[string]*cacheEntry),
		customModules: defaultModules,
	}
}

// cacheEntry is a tuple of the global variables read from a module, and any
// error that may have been observed while parsing the module.
type cacheEntry struct {

	// the mapping of global variable name to Dep
	globals starlark.StringDict

	// any error that was observed while parsing the module
	err error
}

// cachedParser implements Parser by loading Starlark files recursively,
// caching the contents of each file as it goes.
type cachedParser struct {

	// cache is a mapping of module names to a cache entry. The cache allows
	// for the cachedParer to short circuit read operations for files that
	// have already been parsed.
	cache map[string]*cacheEntry

	// root is the root directory with the entrypoint.
	root string

	// reader is a fileReader that will read the dep files.
	reader fileReader

	// customModules is a mapping of Starlark builtin name to Builtin.
	// TODO(nickt): Allow for passing a set of custom modules
	customModules starlark.StringDict
}

func (s cachedParser) Run() error {
	// check for the main entrypoint
	// TODO(nickt): remove the requirement to look for a main file
	main := filepath.Join(s.root, main)
	_, err := oslib.Stat(main)
	if oslib.IsNotExist(err) {
		return fmt.Errorf("main.dep not found")
	}

	load := func(thread *starlark.Thread, moduleName string) (starlark.StringDict, error) {
		var fromPath string
		if thread.CallStackDepth() > 0 {
			fromPath = thread.CallFrame(0).Pos.Filename()
		}
		modulePath, err := s.reader.Resolve(moduleName, fromPath)
		if err != nil {
			return nil, err
		}

		e, ok := s.cache[modulePath]
		if e != nil {
			return e.globals, e.err
		}
		if ok {
			return nil, fmt.Errorf("cycle in load graph")
		}

		moduleSource, err := s.reader.ReadFile(modulePath)
		if err != nil {
			s.cache[modulePath] = &cacheEntry{nil, err}
			return nil, err
		}

		s.cache[modulePath] = nil
		globals, err := starlark.ExecFile(thread, modulePath, moduleSource, s.customModules)
		s.cache[modulePath] = &cacheEntry{globals, err}

		return globals, err
	}

	// recursively load all files
	thread := &starlark.Thread{Load: load}
	_, err = load(thread, main)

	if err != nil {
		return err
	}

	return nil
}

// Deps flattens the deps parsed across all modules and returns a slice of
// pointers to them.
func (s cachedParser) Deps() []*Dep {
	var deps []*Dep
	for _, entry := range s.cache {
		for _, global := range entry.globals {
			if dep, ok := global.(*Dep); ok {
				deps = append(deps, dep)
			}
		}
	}
	return deps
}
