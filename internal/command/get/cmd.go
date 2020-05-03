package get

import (
	"fmt"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

// Command returns a cobra command for `get` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get hog configuration values",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.AddCommand(
		protocol(),
		domain(),
		port(),
	)

	return cmd
}

func protocol() *cobra.Command {
	return &cobra.Command{
		Use:   "protocol",
		Short: "Get protocol value",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			h, err := hog.Get()
			if err != nil {
				utils.PrintError(err)
				return
			}
			fmt.Println(h.Protocol)
		},
	}
}

func domain() *cobra.Command {
	return &cobra.Command{
		Use:   "domain",
		Short: "Get domain value",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			h, err := hog.Get()
			if err != nil {
				utils.PrintError(err)
				return
			}
			fmt.Println(h.Domain)
		},
	}
}

func port() *cobra.Command {
	return &cobra.Command{
		Use:   "port",
		Short: "Get port value",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			h, err := hog.Get()
			if err != nil {
				utils.PrintError(err)
				return
			}
			fmt.Println(h.Port)
		},
	}
}
