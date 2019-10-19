package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildCommit is the commit at which the binary was build. The value is
// provided at build time via -ldflags passed to `go build`.
var BuildCommit string

// NewCommand returns a new command for fetching the version of the binary.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "The version of matryoshka",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(BuildCommit)
			return nil
		},
	}
}
