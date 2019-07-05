package lang

import (
	"math/big"
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestDep_String(t *testing.T) {
	dep := Dep{Name: "foo"}

	expected := "<dep.Dep \"foo\">"
	actual := dep.String()

	if expected != actual {
		t.Errorf("wanted %s, got %s", expected, actual)
	}
}

func TestDep_Type(t *testing.T) {
	dep := Dep{}

	expected := "dep.Dep"
	actual := dep.Type()

	if expected != actual {
		t.Errorf("wanted %s, got %s", expected, actual)
	}
}

func TestDep_Truth(t *testing.T) {
	dep := Dep{}

	if !dep.Truth() {
		t.Error("Truth did not return true")
	}
}

func TestDep_Hash(t *testing.T) {
	dep := Dep{}

	_, err := dep.Hash()

	if err == nil {
		t.Fatal("expected Hash to return an error")
	}

	if !strings.Contains(err.Error(), "unhashable type") {
		t.Errorf("expected error to contain 'unhashable type")
	}
}

func TestAsString(t *testing.T) {
	expected := "foo"
	value := starlark.String(expected)
	actual, err := asString(value)

	if err != nil {
		t.Errorf("could not parse starlark string")
	}

	if expected != actual {
		t.Errorf("wanted %s, got %s", expected, actual)
	}
}

func TestAsString_NotString(t *testing.T) {
	values := []starlark.Value{
		starlark.MakeInt(42),
		starlark.MakeBigInt(big.NewInt(42)),
		starlark.NewDict(42),
		starlark.NewList([]starlark.Value{}),
	}

	for _, value := range values {
		_, err := asString(value)
		if err == nil {
			t.Errorf("expected error parsing value %v to a string", value)
		}
	}
}

func TestAsRequirementsList(t *testing.T) {
	depItems := []*Dep{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}
	var values []starlark.Value
	for _, item := range depItems {
		values = append(values, item)
	}

	list, err := asRequirementsList(starlark.NewList(values))
	if err != nil {
		t.Errorf("error parsing requirements list: %s", err)
	}

	for i, item := range list {
		expected := depItems[i]
		if expected != item {
			t.Errorf("wanted %s, got %s", expected, item)
		}
	}
}

func TestAsRequirementsList_NotDep(t *testing.T) {
	lists := []*starlark.List{
		starlark.NewList([]starlark.Value{starlark.String("foo")}),
		starlark.NewList([]starlark.Value{starlark.MakeInt(42)}),
		starlark.NewList([]starlark.Value{starlark.NewDict(42)}),
	}

	for _, list := range lists {
		_, err := asRequirementsList(list)
		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "is not a dep") {
			t.Error("expected error message to contain 'is not a dep'")
		}
	}
}

func TestAsRequirementsList_NotList(t *testing.T) {
	values := []starlark.Value{
		starlark.String("foo"),
		starlark.NewDict(42),
		starlark.MakeInt(42),
	}

	for _, value := range values {
		_, err := asRequirementsList(value)
		if err == nil {
			t.Fatalf("expected error parsing %s as a list", value)
		}

		if !strings.Contains(err.Error(), "is not a list") {
			t.Errorf("expected error message to contain 'is not a list'")
		}
	}
}

func TestAsCommandList(t *testing.T) {
	cmdItems := []ShellCmd{
		{"foo"},
		{"bar"},
		{"baz"},
	}
	var values []starlark.Value
	for _, item := range cmdItems {
		values = append(values, item)
	}

	list, err := asCommandList(starlark.NewList(values))
	if err != nil {
		t.Errorf("error parsing requirements list: %s", err)
	}

	for i, item := range list {
		expected := cmdItems[i].Command
		if expected != item.Command {
			t.Errorf("wanted %s, got %s", expected, item)
		}
	}
}

func TestAsCommandList_NotDep(t *testing.T) {
	lists := []*starlark.List{
		starlark.NewList([]starlark.Value{starlark.String("foo")}),
		starlark.NewList([]starlark.Value{starlark.MakeInt(42)}),
		starlark.NewList([]starlark.Value{starlark.NewDict(42)}),
	}

	for _, list := range lists {
		_, err := asCommandList(list)
		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "is not a command") {
			t.Error("expected error message to contain 'is not a command'")
		}
	}
}

func TestAsCommandList_NotList(t *testing.T) {
	values := []starlark.Value{
		starlark.String("foo"),
		starlark.NewDict(42),
		starlark.MakeInt(42),
	}

	for _, value := range values {
		_, err := asCommandList(value)
		if err == nil {
			t.Fatalf("expected error parsing %s as a list", value)
		}

		if !strings.Contains(err.Error(), "is not a list") {
			t.Errorf("expected error message to contain 'is not a list'")
		}
	}
}
