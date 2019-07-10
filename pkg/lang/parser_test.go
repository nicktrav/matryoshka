package lang

import (
	"testing"

	"go.starlark.net/starlark"
)

const (
	noMain     = "./testcases/nomain"
	simpleMain = "./testcases/simple_main"
	multiFile  = "./testcases/multi_file"
)

func TestParser_NoMain(t *testing.T) {
	parser := NewParser(noMain)

	err := parser.Run()
	if err == nil {
		t.Fatalf("wanted an erorr")
	}

	if err.Error() != "main.dep not found" {
		t.Errorf("wanted error 'main.dep not found'")
	}
}

func TestParser_SimpleMain(t *testing.T) {
	parser := NewParser(simpleMain)

	err := parser.Run()
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	modules := parser.Modules()
	if len(modules) != 1 {
		t.Errorf("wanted 1 file , got %d", len(modules))
	}

	main := modules[0]
	if len(main.Keys()) != 2 {
		t.Errorf("wanted 2 variables in main.dep, got %d", len(main.Keys()))
	}

	// Check the "all" Dep

	all, ok := main["all"]
	if !ok {
		t.Errorf("main.dep does to have Dep 'all'")
	}

	dep, ok := all.(*Dep)
	if !ok {
		t.Errorf("'all' was not a Dep")
	}

	if dep.Name != "all" {
		t.Errorf("name of Dep 'all' was not 'all'; got %s", dep.Name)
	}

	if len(dep.Requirements) != 1 {
		t.Errorf("wanted Dep 'all' to have 1 requirement; got %d", len(dep.Requirements))
	}

	if dep.Requirements[0].Name != "foo" {
		t.Errorf("wanted Dep 'all' to have requirement 'foo'; got '%s'", dep.Requirements[0])
	}

	if len(dep.MetCommands) != 1 {
		t.Errorf("wanted Dep 'all' to have 1 'met' command; got %d", len(dep.MetCommands))
	}

	want := "echo 'Hello, world!'"
	got := dep.MetCommands[0].Command
	if dep.MetCommands[0].Command != "echo 'Hello, world!'" {
		t.Errorf("wanted Dep 'all' 'met' command to be '%s'; got '%s'", want, got)
	}

	if len(dep.MeetCommands) != 2 {
		t.Errorf("wanted Dep 'all' to have 2 'meet' command; got %d", len(dep.MeetCommands))
	}

	want = "echo 'Hello, indeed!'"
	got = dep.MeetCommands[0].Command
	if want != got {
		t.Errorf("wanted Dep 'all' 'meet' command #1 to be '%s'; got '%s'", want, got)
	}

	want = "echo 'Hello, again!'"
	got = dep.MeetCommands[1].Command
	if want != got {
		t.Errorf("wanted Dep 'all' 'meet' command #2 to be '%s'; got '%s'", want, got)
	}

	// Check the "foo" Dep

	foo, ok := main["foo"]
	if !ok {
		t.Errorf("wanted main.dep to have dep 'foo'")
	}

	dep, ok = foo.(*Dep)
	if !ok {
		t.Errorf("wanted 'foo' to be a Dep")
	}
	if dep.Name != "foo" {
		t.Errorf("wanted name of Dep 'foo' to be 'foo'; got %s", dep.Name)
	}

	if len(dep.Requirements) != 0 {
		t.Errorf("wanted Dep 'foo' to have 0 requirement; got %d", len(dep.Requirements))
	}

	if len(dep.MetCommands) != 0 {
		t.Errorf("wanted Dep 'foo' to have 0 'met' command; got %d", len(dep.MetCommands))
	}

	if len(dep.MeetCommands) != 0 {
		t.Errorf("wanted Dep 'foo' to have 0 'meet' command; got %d", len(dep.MeetCommands))
	}
}

func TestParser_MultiFile(t *testing.T) {
	parser := NewParser(multiFile)

	err := parser.Run()
	if err != nil {
		t.Errorf("parser.Run: %s", err)
	}

	modules := parser.Modules()
	if len(modules) != 4 {
		t.Errorf("wanted 4 files, got %d", len(modules))
	}

	depMap := toMap(modules)
	wantedGlobals := []string{"all", "foo", "bar", "baz", "bam"}
	for _, want := range wantedGlobals {
		if _, contains := depMap[want]; !contains {
			t.Errorf("wanted module %s in module %+v", wantedGlobals, modules)
		}
	}

	// all depends on foo, bar and baz

	dep := depMap["all"]
	if !contains("foo", dep.Requirements) {
		t.Errorf("wanted Dep 'all' to require 'foo'")
	}

	if !contains("bar", dep.Requirements) {
		t.Errorf("wanted Dep 'all' to require 'bar'")
	}

	if !contains("baz", dep.Requirements) {
		t.Errorf("wanted Dep 'all' to require 'baz'")
	}

	// foo depends on bam

	dep = depMap["foo"]
	if !contains("bam", dep.Requirements) {
		t.Errorf("wanted Dep 'foo' to require 'bam'")
	}

	// bar depends on baz

	dep = depMap["bar"]
	if !contains("baz", dep.Requirements) {
		t.Errorf("wanted Dep 'bar' to require 'baz'")
	}

	// baz depends on nothing

	dep = depMap["baz"]
	if len(dep.Requirements) != 0 {
		t.Errorf("wanted Dep 'baz' to depend on nothing; got %s", dep.Requirements)
	}

	// bam depends on baz

	dep = depMap["bam"]
	if !contains("baz", dep.Requirements) {
		t.Errorf("wanted Dep 'bam' to require 'baz'")
	}
}

func toMap(modules []starlark.StringDict) map[string]*Dep {
	depMap := make(map[string]*Dep)

	for _, module := range modules {
		for key, val := range module {
			dep, ok := val.(*Dep)
			if !ok {
				continue
			}
			depMap[key] = dep
		}
	}

	return depMap
}

func contains(want string, requirements []*Dep) bool {
	for _, req := range requirements {
		if req.Name == want {
			return true
		}
	}
	return false
}
