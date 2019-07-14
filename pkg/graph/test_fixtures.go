package graph

import "github.com/nicktrav/matryoshka/pkg/lang"

var fooRawDep = &lang.Dep{
	Name:         "foo",
	Requirements: []*lang.Dep{barRawDep, bamRawDep},
	MetCommands:  []*lang.ShellCmd{},
	MeetCommands: []*lang.ShellCmd{},
}

var barRawDep = &lang.Dep{
	Name:         "bar",
	Requirements: []*lang.Dep{bazRawDep, boomRawDep},
	MetCommands:  []*lang.ShellCmd{},
	MeetCommands: []*lang.ShellCmd{},
}

var bazRawDep = &lang.Dep{
	Name:         "baz",
	Requirements: []*lang.Dep{},
	MetCommands:  []*lang.ShellCmd{},
	MeetCommands: []*lang.ShellCmd{},
}

var bamRawDep = &lang.Dep{
	Name:         "bam",
	Requirements: []*lang.Dep{boomRawDep},
	MetCommands:  []*lang.ShellCmd{},
	MeetCommands: []*lang.ShellCmd{},
}

var boomRawDep = &lang.Dep{
	Name:         "boom",
	Requirements: []*lang.Dep{},
	MetCommands:  []*lang.ShellCmd{},
	MeetCommands: []*lang.ShellCmd{},
}
