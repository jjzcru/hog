package remove

import (
"fmt"
"github.com/spf13/cobra"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove file reference from the service",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Remove command")
		},
	}

	return cmd
}

