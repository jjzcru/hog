package add
import (
	"fmt"
	"github.com/spf13/cobra"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Make file/s accessible to the service",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add command")
		},
	}

	return cmd
}
